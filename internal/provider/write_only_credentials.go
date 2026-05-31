package provider

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/crypto/argon2"
)

const (
	writeOnlyCredentialVerifiersPrivateKey = "write_only_credential_verifiers"

	writeOnlyCredentialArgon2IDMemory      = 64 * 1024
	writeOnlyCredentialArgon2IDTime        = 1
	writeOnlyCredentialArgon2IDParallelism = 4
	writeOnlyCredentialSaltLength          = 16
	writeOnlyCredentialKeyLength           = 32
	writeOnlyCredentialParameterPartCount  = 3
)

type privateStateReader interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
}

type privateStateWriter interface {
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

func writeOnlyStringConfigured(value types.String) bool {
	return !value.IsNull() && !value.IsUnknown()
}

func validateStringCredential(diags *diag.Diagnostics, legacy, writeOnly types.String, legacyPath, writeOnlyPath path.Path) {
	_, _, credentialDiags := resolveStringCredential(legacy, writeOnly, legacyPath, writeOnlyPath, false)
	diags.Append(credentialDiags...)
}

func resolveStringCredential(legacy, writeOnly types.String, legacyPath, writeOnlyPath path.Path, required bool) (types.String, bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	legacyConfigured := !legacy.IsNull() && !legacy.IsUnknown()
	writeOnlyConfigured := writeOnlyStringConfigured(writeOnly)

	if legacyConfigured && writeOnlyConfigured {
		diags.AddError(
			"Conflicting credential arguments",
			"Only one of "+legacyPath.String()+" and "+writeOnlyPath.String()+" can be configured.",
		)

		return types.StringNull(), false, diags
	}

	if writeOnlyConfigured {
		return types.StringValue(writeOnly.ValueString()), true, diags
	}

	if legacyConfigured {
		return legacy, false, diags
	}

	if required {
		diags.AddError(
			"Missing credential argument",
			"One of "+legacyPath.String()+" or "+writeOnlyPath.String()+" must be configured.",
		)
	}

	return types.StringNull(), false, diags
}

type writeOnlyCredentialValues map[string]types.String

func (v writeOnlyCredentialValues) Add(argumentPath path.Path, value types.String) {
	v[argumentPath.String()] = value
}

func writeOnlyCredentialPreimage(argumentPath string, value types.String) []byte {
	return []byte(argumentPath + "\x00" + value.ValueString())
}

func writeOnlyCredentialVerifier(argumentPath string, value types.String) (string, error) {
	salt := make([]byte, writeOnlyCredentialSaltLength)

	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("generate write-only credential verifier salt: %w", err)
	}

	key := argon2.IDKey(
		writeOnlyCredentialPreimage(argumentPath, value),
		salt,
		writeOnlyCredentialArgon2IDTime,
		writeOnlyCredentialArgon2IDMemory,
		writeOnlyCredentialArgon2IDParallelism,
		writeOnlyCredentialKeyLength,
	)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		writeOnlyCredentialArgon2IDMemory,
		writeOnlyCredentialArgon2IDTime,
		writeOnlyCredentialArgon2IDParallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func writeOnlyCredentialVerifierMatches(argumentPath string, value types.String, verifier string) bool {
	parts := strings.Split(verifier, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false
	}

	version, hasVersion := strings.CutPrefix(parts[2], "v=")
	if !hasVersion {
		return false
	}

	parsedVersion, err := strconv.Atoi(version)
	if err != nil || parsedVersion != argon2.Version {
		return false
	}

	parameters, hasMemoryParameter := strings.CutPrefix(parts[3], "m=")
	if !hasMemoryParameter {
		return false
	}

	parameterParts := strings.Split(parameters, ",")
	if len(parameterParts) != writeOnlyCredentialParameterPartCount {
		return false
	}

	memory, err := strconv.ParseUint(parameterParts[0], 10, 32)
	if err != nil {
		return false
	}

	timePart, hasTimeParameter := strings.CutPrefix(parameterParts[1], "t=")
	if !hasTimeParameter {
		return false
	}

	time, err := strconv.ParseUint(timePart, 10, 32)
	if err != nil {
		return false
	}

	parallelismPart, hasParallelismParameter := strings.CutPrefix(parameterParts[2], "p=")
	if !hasParallelismParameter {
		return false
	}

	parallelism, err := strconv.ParseUint(parallelismPart, 10, 8)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	expectedKey, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	if len(expectedKey) != writeOnlyCredentialKeyLength {
		return false
	}

	actualKey := argon2.IDKey(
		writeOnlyCredentialPreimage(argumentPath, value),
		salt,
		uint32(time),
		uint32(memory),
		uint8(parallelism),
		writeOnlyCredentialKeyLength,
	)

	return subtle.ConstantTimeCompare(actualKey, expectedKey) == 1
}

func readWriteOnlyCredentialVerifiers(ctx context.Context, private privateStateReader) (map[string]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if private == nil {
		return map[string]string{}, diags
	}

	data, getDiags := private.GetKey(ctx, writeOnlyCredentialVerifiersPrivateKey)
	diags.Append(getDiags...)

	if diags.HasError() || len(data) == 0 {
		return map[string]string{}, diags
	}

	verifiers := map[string]string{}

	err := json.Unmarshal(data, &verifiers)
	if err != nil {
		diags.AddError("Invalid write-only credential verifier private state", err.Error())

		return map[string]string{}, diags
	}

	return verifiers, diags
}

func writeWriteOnlyCredentialVerifiers(ctx context.Context, private privateStateWriter, values writeOnlyCredentialValues) diag.Diagnostics {
	var diags diag.Diagnostics

	if private == nil {
		return diags
	}

	if len(values) == 0 {
		diags.Append(private.SetKey(ctx, writeOnlyCredentialVerifiersPrivateKey, nil)...)

		return diags
	}

	verifiers := make(map[string]string, len(values))
	for argumentPath, value := range values {
		verifier, err := writeOnlyCredentialVerifier(argumentPath, value)
		if err != nil {
			diags.AddError("Failed to create write-only credential verifier", err.Error())

			return diags
		}

		verifiers[argumentPath] = verifier
	}

	data, err := json.Marshal(verifiers)
	if err != nil {
		diags.AddError("Invalid write-only credential verifiers", err.Error())

		return diags
	}

	diags.Append(private.SetKey(ctx, writeOnlyCredentialVerifiersPrivateKey, data)...)

	return diags
}

func writeOnlyCredentialVerifiersChanged(ctx context.Context, private privateStateReader, configured writeOnlyCredentialValues) (bool, diag.Diagnostics) {
	stored, diags := readWriteOnlyCredentialVerifiers(ctx, private)
	if diags.HasError() {
		return false, diags
	}

	if len(stored) != len(configured) {
		return true, diags
	}

	for argumentPath, value := range configured {
		if !writeOnlyCredentialVerifierMatches(argumentPath, value, stored[argumentPath]) {
			return true, diags
		}
	}

	return false, diags
}

func markWriteOnlyCredentialChange(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, values writeOnlyCredentialValues, connectionDetails attr.Value) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	changed, verifierDiags := writeOnlyCredentialVerifiersChanged(ctx, req.Private, values)
	resp.Diagnostics.Append(verifierDiags...)

	if resp.Diagnostics.HasError() || !changed {
		return
	}

	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("connection_details"), connectionDetails)...)
}

type writeOnlyCredentialModelResolver[T any] func(plan, config T) (T, writeOnlyCredentialValues, diag.Diagnostics)

func modifyPlanForWriteOnlyCredentialChange[T any](
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
	resolver writeOnlyCredentialModelResolver[T],
	connectionDetails attr.Value,
) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var plan, config T

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, values, modelDiags := resolver(plan, config)
	resp.Diagnostics.Append(modelDiags...)

	markWriteOnlyCredentialChange(ctx, req, resp, values, connectionDetails)
}
