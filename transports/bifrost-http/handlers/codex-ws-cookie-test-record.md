# FWD-67：原生 WS 握手 Cookie 测试结果

需求、调研、设计：用户收到 WS Set-Cookie 缺口与调整握手顺序、清理/续聊验收建议后授权继续开发和 push，之后统一评审。本项基底2990e6dda；复用已有 fasthttp/websocket，无新依赖。成功上游握手在客户端101之前完成，只转发原始Set-Cookie，安全握手头仍由各跳库生成。预留连接名额后拨号，失败释放；成功连接交给同一session，各轮业务仍经治理。未合并、未部署、未真机、无新抓包。

执行：GPT/Codex，2026-10-04 AEDT；Windows本机Go，显式本地go.work引用仓库内模块。假公司token、假虚拟key、loopback HTTP/WS上游，无真实凭据、无付费API。测试编号FWD-67-A2~。go.work只作本机环境，不提交；不存在则在根目录执行：

```powershell
go work init ./core ./framework ./transports ./plugins/governance ./plugins/logging ./plugins/telemetry ./plugins/maxim ./plugins/otel ./plugins/modelcatalogresolver ./plugins/jsonparser ./plugins/routing ./plugins/semanticcache ./plugins/mocker ./plugins/compat ./plugins/prompts
$env:GOWORK = (Join-Path (Get-Location) go.work)
go test -json ./transports/bifrost-http/handlers ./transports/bifrost-http/integrations ./transports/bifrost-http/websocket -run 'TestChatGPT|TestCodexWSIdentity|TestSession' -count=1
go test -json ./transports/bifrost-http/handlers ./transports/bifrost-http/integrations ./transports/bifrost-http/websocket -count=1
```

| 验收动作 | 脚本/用例 | 结果 |
|---|---|---|
| 同名不同路径、Expires逗号、多条Cookie原值顺序；两轮复用，第3轮限额 | TestChatGPTWSHandshakeCookies | 通过 |
| 上游401→通用502，无Cookie外露；客户端升级拒绝→403，上游和名额释放 | TestChatGPTWSHandshakeFailureCleanup | 2子场景通过 |
| 无Cookie不新增；第二连接429且不拨上游；无消息退出清理 | TestChatGPTWSHandshakeAdmission | 通过 |
| 预留限额、幂等释放、完成、关闭后不可恢复旧预留 | TestSessionReservationLimitAndCleanup | 通过 |
| 既有身份映射、多轮原字节、每轮限额、断连与上下文 | TestChatGPT / TestCodexWSIdentity / TestSession | 定向通过 |

修复前red：TestChatGPTWSHandshakeCookies失败，cookies=[]，预期两条；未发消息时旧实现尚未拨上游，因此等待上游关闭也超时，此项不是已连接上游泄漏证据。修复后green：该用例PASS 0.196s。最初go.work未实际生成导致go.mod/文件缺失，两次均未执行断言；提升权限创建后正常。离线harness结构检查初轮EPERM，提升权限后成功；生成639请求，过滤Codex WebSocket得到1个已有HTTP入口guard。Newman不能执行WS握手Cookie验收，collection只更新覆盖说明，不能冒称新增端到端Cookie harness用例；真实本地WS回归覆盖，真机留作统一评审后复验。未运行真实provider harness或make test-core付费接口。

定向真实计数：{"pass": 42, "fail": 0, "skip": 0}；修复版三包：{"pass": 1314, "fail": 104, "skip": 0}；确切基底2990e6dda三包：{"pass": 1308, "fail": 104, "skip": 0}。计数含子测试，三包结束事件齐全。104个既有失败名称集合完全相同；全量非全绿，不报整个FWD-67已完成。

## 统计

```json
{
  "source": "C:\\Users\\MI\\.codex\\sessions\\2026\\10\\04\\rollout-2026-10-04T13-33-06-01a104c2-0331-71e1-b62b-3a085d8fc2bd.jsonl",
  "start": "2026-10-04T03:53:32.751Z",
  "archive_cutoff": "2026-10-04T15:08:40.837774+11:00",
  "model": [
    "gpt-6.1-sol"
  ],
  "token_delta": {
    "input_tokens": 2931646,
    "cached_input_tokens": 2823552,
    "cache_write_input_tokens": 0,
    "output_tokens": 14803,
    "reasoning_output_tokens": 2835,
    "total_tokens": 2946449
  },
  "scope": "本批授权至本项封存的开发session累计差值，含缓存；不计审批会话及随后提交/push收尾；整个FWD-67跨轮累计未核实",
  "elapsed_seconds": 908
}
```

## 全量失败名称（修复版=本次基底）

```text
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestClaudeMarketplaceGitRepoContainsMarketplaceAndCloneablePlugin
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestCleanupOrphanSkillFilesDeletesDBFallbackBlobs
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestCleanupOrphanSkillFilesDeletesOnlyUnreferencedUploadObjects
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestComplexityAnalyzerConfigGetDegradesOnUnreadableConfig
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestComplexityAnalyzerConfigGetReturnsDefaultsWhenUnset
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestComplexityAnalyzerConfigPutPersistsAndReloads
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestComplexityAnalyzerConfigPutRejectsInvalidPayloads
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestComplexityAnalyzerConfigResetPersistsDefaultsAndReloads
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestComplexityAnalyzerConfigResetReportsReloadFailure
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestEnrichListModelsResponse_MarksDeprecatedPricingRows
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestGetModelParameters_ResolvesQualifiedAndBareIDs
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestGetVirtualKeyQuota_EndToEndWithRealStore
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestGetVirtualKeyQuota_ExpiredGraceValueUnauthorized
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestGetVirtualKeyQuota_GraceValueWithRealStore
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestGetVirtualKeyQuota_WindowClampedToBudgetCreation
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleAuthorize
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleAuthorize/happy_path_redirects_to_the_consent_page
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleAuthorize/loopback_redirect_matches_on_any_port
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleAuthorize/mismatched_resource_redirects_with_error
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleAuthorize/non-S256_challenge_method_redirects_with_error
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleAuthorize/non-code_response_type_redirects_with_error
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleAuthorize/omitted_resource_defaults_to_canonical_and_proceeds
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleAuthorize/scope_exceeds_registered_redirects_with_error
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleAuthorize/unknown_client_id_is_rejected
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleAuthorize/unregistered_redirect_uri_is_rejected
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleAuthorize_AuthCodeTTLResolution
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleAuthorize_AuthCodeTTLResolution/above_cap_is_clamped
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleAuthorize_AuthCodeTTLResolution/in-range_used_verbatim
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleAuthorize_AuthCodeTTLResolution/zero_falls_back_to_default
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleRegister_DCR
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleRegister_DCR/malformed_JSON_is_rejected
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleRegister_DCR/missing_redirect_uris_is_rejected
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleRegister_DCR/non-public_auth_method_is_rejected
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleRegister_DCR/unsupported_grant_type_is_rejected
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleRegister_DCR/unsupported_response_type_is_rejected
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleRegister_DCR/valid_registration_returns_201_with_defaults
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_AuthorizationCode
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_AuthorizationCode/PKCE_mismatch_is_rejected
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_AuthorizationCode/client_id_mismatch_is_rejected
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_AuthorizationCode/code_is_single-use
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_AuthorizationCode/expired_code_is_rejected
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_AuthorizationCode/happy_path_issues_a_verifiable_token_pair
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_AuthorizationCode/missing_redirect_uri_is_rejected
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_AuthorizationCode/missing_required_fields_are_rejected
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_AuthorizationCode/unsupported_grant_type_is_rejected
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_RefreshRotationAndReplay
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_RefreshRotationAndReplay/client_id_mismatch_is_rejected
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_RefreshRotationAndReplay/missing_fields_are_rejected
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_RefreshRotationAndReplay/replaying_a_rotated_token_revokes_the_whole_family
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_RefreshRotationAndReplay/rotation_issues_a_new_pair_and_carries_the_family
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_RefreshUserLiveness
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_RefreshUserLiveness/user_refresh_rejected_when_the_user_is_no_longer_active
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_RefreshUserLiveness/user_refresh_succeeds_when_the_user_is_active
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_RefreshVKIdentityDisabled
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_RefreshVKIdentityDisabled/vk_refresh_not_blocked_by_flag_when_user_mode_unavailable
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestHandleToken_RefreshVKIdentityDisabled/vk_refresh_rejected_when_disabled_and_user_mode_available
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestListBaseModels_IncludesDeprecatedPricingRows
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestListModelDetails_AppliesGlobalOverrideWithoutMutatingBase
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestListModelDetails_AppliesProviderScopedOverride
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestListModelDetails_IncludesPricing
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestListModelDetails_NoOverridesOmitsNewFields
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestListModelDetails_OverrideIndexIsDeduplicated
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestListModelDetails_OverrideWithoutBaseCatalogRow
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestListModelDetails_PatchToSameValueIsNotMarkedOverridden
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestListModelDetails_ResolvesCatalogPricing
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestListModelDetails_ResolvesCatalogPricing/Azure_alias
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestListModelDetails_ResolvesCatalogPricing/Azure_alias_falls_back_to_alias_key
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestListModelDetails_ResolvesCatalogPricing/Azure_alias_with_empty_model_name
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestListModelDetails_ResolvesCatalogPricing/Together_alias
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestListModelDetails_ResolvesCatalogPricing/Together_catalog_provider
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestListModelDetails_VirtualKeyScopedOverrideIsInformationalOnly
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestListModels_MarksDeprecatedModelsWithoutFiltering
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestProviderGovernance_DecodesEncodedProviderName
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestProviderGovernance_MalformedEncodingReturns400
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestProviderGovernance_UnknownProviderStill404
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestRetryComplexitySemanticWarmup
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestRetryComplexitySemanticWarmup/accepts_a_failed_warmup
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestRetryComplexitySemanticWarmup/rejects_a_non-failed_warmup
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestSkillsServingGenericFileDownloadDecodesEncodedPathParams
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestTeam_DecodesEncodedTeamID
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestTeam_MalformedEncodingReturns400
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestUpdateConfig_PersistsVKRotationCooldown
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestUpdateConfig_RejectsAuthCodeTTLAboveMax
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestUpdateConfig_RejectsAuthCodeTTLAboveMax/both
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestUpdateConfig_RejectsAuthCodeTTLAboveMax/headers
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestUpdateConfig_RejectsAuthCodeTTLAboveMax/oauth
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestUpdatePricingOverride_ReplacesFullBody
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestUpdateProviderGovernance_AdoptsReconciledBudgetsNotStaleOnes
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestVirtualKeyBudgetOverrideLifecycle
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestVirtualKeyBudgetOverrideRejectsDirectMirrorBudget
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestWebhookHandlerCreateAndGet
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestWebhookHandlerCreateValidation
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestWebhookHandlerDelete
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestWebhookHandlerDuplicateName
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestWebhookHandlerHeaders
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestWebhookHandlerListDeliveries
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestWebhookHandlerListFilters
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestWebhookHandlerPerEndpointTuning
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestWebhookHandlerRedeliver
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestWebhookHandlerRotateSecret
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestWebhookHandlerRouteRegistration
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestWebhookHandlerTestDelivery
github.com/maximhq/bifrost/transports/bifrost-http/handlers/TestWebhookHandlerUpdate
github.com/maximhq/bifrost/transports/bifrost-http/integrations/TestGenericRouter_MarkDeprecatedListModelsResponseUsesCatalog
```

## 可取得真实输出


### 定向

```text
=== RUN   TestChatGPTPassthroughRouterRegistersCodexResponsesPost
--- PASS: TestChatGPTPassthroughRouterRegistersCodexResponsesPost (0.00s)
=== RUN   TestChatGPTBackgroundRoutes
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/codex/models?client_version=0.159.3
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/codex/models?client_version=0.159.3 (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/HEAD/backend-api/codex/models
--- PASS: TestChatGPTBackgroundRoutes/HEAD/backend-api/codex/models (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/ps/plugins/list
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/ps/plugins/list (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/POST/backend-api/ps/plugins/install
--- PASS: TestChatGPTBackgroundRoutes/POST/backend-api/ps/plugins/install (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/ps/mcp
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/ps/mcp (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/PUT/backend-api/ps/mcp/config
--- PASS: TestChatGPTBackgroundRoutes/PUT/backend-api/ps/mcp/config (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/wham/usage
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/wham/usage (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/plugins/featured
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/plugins/featured (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/POST/backend-api/codex/analytics-events/events
--- PASS: TestChatGPTBackgroundRoutes/POST/backend-api/codex/analytics-events/events (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/codex/analytics-events/events
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/codex/analytics-events/events (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/POST/backend-api/codex/models
--- PASS: TestChatGPTBackgroundRoutes/POST/backend-api/codex/models (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/DELETE/backend-api/accounts
--- PASS: TestChatGPTBackgroundRoutes/DELETE/backend-api/accounts (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/ps/mcp-lookalike
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/ps/mcp-lookalike (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/ps/plugins-other/list
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/ps/plugins-other/list (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/TRACE/backend-api/wham/usage
--- PASS: TestChatGPTBackgroundRoutes/TRACE/backend-api/wham/usage (0.00s)
--- PASS: TestChatGPTBackgroundRoutes (0.00s)
=== RUN   TestChatGPTBackgroundWireAndHooks
--- PASS: TestChatGPTBackgroundWireAndHooks (0.02s)
=== RUN   TestChatGPTUpstreamURLDefaultsToChatGPT
--- PASS: TestChatGPTUpstreamURLDefaultsToChatGPT (0.00s)
=== RUN   TestChatGPTUpstreamURLHonoursDiagnosticOverride
--- PASS: TestChatGPTUpstreamURLHonoursDiagnosticOverride (0.00s)
PASS
=== RUN   TestSessionManagerCreateAndGet
--- PASS: TestSessionManagerCreateAndGet (0.00s)
=== RUN   TestSessionManagerConnectionLimit
--- PASS: TestSessionManagerConnectionLimit (0.00s)
=== RUN   TestSessionReservationLimitAndCleanup
--- PASS: TestSessionReservationLimitAndCleanup (0.00s)
=== RUN   TestSessionManagerRemove
--- PASS: TestSessionManagerRemove (0.00s)
=== RUN   TestSessionLastResponseID
--- PASS: TestSessionLastResponseID (0.00s)
=== RUN   TestSessionManagerCloseAll
--- PASS: TestSessionManagerCloseAll (0.00s)
=== RUN   TestSessionRealtimeState
--- PASS: TestSessionRealtimeState (0.00s)
=== RUN   TestSessionRecordRealtimeInputUpdatesPendingItemAndIgnoresConsumedLateUpdate
--- PASS: TestSessionRecordRealtimeInputUpdatesPendingItemAndIgnoresConsumedLateUpdate (0.00s)
PASS
ok  	github.com/maximhq/bifrost/transports/bifrost-http/integrations	0.169s
ok  	github.com/maximhq/bifrost/transports/bifrost-http/websocket	0.171s
=== RUN   TestCodexWSIdentityPreservesBusinessBytes
--- PASS: TestCodexWSIdentityPreservesBusinessBytes (0.00s)
=== RUN   TestCodexWSIdentityRejectsInvalidContext
--- PASS: TestCodexWSIdentityRejectsInvalidContext (0.00s)
=== RUN   TestSessionLoginInvalidPayloadDoesNotExposeDecoderDetails
--- PASS: TestSessionLoginInvalidPayloadDoesNotExposeDecoderDetails (0.00s)
=== RUN   TestChatGPTWSHandshakeCookies
--- PASS: TestChatGPTWSHandshakeCookies (0.10s)
=== RUN   TestChatGPTWSHandshakeFailureCleanup
=== RUN   TestChatGPTWSHandshakeFailureCleanup/false
--- PASS: TestChatGPTWSHandshakeFailureCleanup/false (0.00s)
=== RUN   TestChatGPTWSHandshakeFailureCleanup/true
--- PASS: TestChatGPTWSHandshakeFailureCleanup/true (0.00s)
--- PASS: TestChatGPTWSHandshakeFailureCleanup (0.00s)
=== RUN   TestChatGPTWSHandshakeAdmission
--- PASS: TestChatGPTWSHandshakeAdmission (0.00s)
=== RUN   TestChatGPTWSMultiTurnRawAndQuota
=== RUN   TestChatGPTWSMultiTurnRawAndQuota/false
--- PASS: TestChatGPTWSMultiTurnRawAndQuota/false (0.00s)
=== RUN   TestChatGPTWSMultiTurnRawAndQuota/true
--- PASS: TestChatGPTWSMultiTurnRawAndQuota/true (0.00s)
--- PASS: TestChatGPTWSMultiTurnRawAndQuota (0.00s)
=== RUN   TestChatGPTWSTargetRejectsPlaintextRemote
--- PASS: TestChatGPTWSTargetRejectsPlaintextRemote (0.00s)
=== RUN   TestChatGPTWSRouteRequiresUpgradeAndCredentials
--- PASS: TestChatGPTWSRouteRequiresUpgradeAndCredentials (0.00s)
=== RUN   TestChatGPTWSDisconnectDuringTurnFinalizesHooks
--- PASS: TestChatGPTWSDisconnectDuringTurnFinalizesHooks (0.00s)
PASS
ok  	github.com/maximhq/bifrost/transports/bifrost-http/handlers	0.351s

```

### 修复版全量

```text
=== RUN   TestPoolGetAndReturn
--- PASS: TestPoolGetAndReturn (0.02s)
=== RUN   TestPoolMaxIdlePerKey
--- PASS: TestPoolMaxIdlePerKey (0.00s)
=== RUN   TestPoolClose
--- PASS: TestPoolClose (0.00s)
=== RUN   TestPoolDialIncludesHandshakeDetails
--- PASS: TestPoolDialIncludesHandshakeDetails (0.00s)
=== RUN   TestDialUpstreamCloseDoesNotAffectPoolCapacityCounters
--- PASS: TestDialUpstreamCloseDoesNotAffectPoolCapacityCounters (0.00s)
=== RUN   TestPoolExpiredConnection
=== RUN   TestAnthropicRawStreamTextCodecRewritesOnlyTextDelta
--- PASS: TestAnthropicRawStreamTextCodecRewritesOnlyTextDelta (0.00s)
=== RUN   TestAnthropicRawStreamTextCodecIgnoresNonTextEvents
--- PASS: TestAnthropicRawStreamTextCodecIgnoresNonTextEvents (0.00s)
=== RUN   TestAnthropicRawStreamTextCodecRejectsMalformedEligibleEvents
--- PASS: TestAnthropicRawStreamTextCodecRejectsMalformedEligibleEvents (0.00s)
=== RUN   TestRewriteAnthropicRawRequestBodyRedactsOnlyContentFields
--- PASS: TestRewriteAnthropicRawRequestBodyRedactsOnlyContentFields (0.00s)
=== RUN   TestRewriteAnthropicRawRequestBodyRejectsUnmappedLiteral
--- PASS: TestRewriteAnthropicRawRequestBodyRejectsUnmappedLiteral (0.00s)
=== RUN   TestRewriteAnthropicRawRequestBodyRejectsMalformedJSON
--- PASS: TestRewriteAnthropicRawRequestBodyRejectsMalformedJSON (0.00s)
=== RUN   TestRewriteAnthropicRawRequestBodyRejectsDuplicateKeys
--- PASS: TestRewriteAnthropicRawRequestBodyRejectsDuplicateKeys (0.00s)
=== RUN   TestRewriteAnthropicRawRequestBodyTransformsTargetsDuplicateText
--- PASS: TestRewriteAnthropicRawRequestBodyTransformsTargetsDuplicateText (0.00s)
=== RUN   TestRewriteAnthropicRawRequestBodyTransformsRejectsOriginalMismatch
--- PASS: TestRewriteAnthropicRawRequestBodyTransformsRejectsOriginalMismatch (0.00s)
=== RUN   TestRewriteAnthropicRawRequestBodyTransformsPreservesHistory
--- PASS: TestRewriteAnthropicRawRequestBodyTransformsPreservesHistory (0.00s)
=== RUN   TestRewriteAnthropicRawResponseTransformsTargetsDuplicateText
--- PASS: TestRewriteAnthropicRawResponseTransformsTargetsDuplicateText (0.00s)
=== RUN   TestMustConvertInPassthrough
=== RUN   TestMustConvertInPassthrough/added_message
--- PASS: TestMustConvertInPassthrough/added_message (0.00s)
=== RUN   TestMustConvertInPassthrough/added_advisor
--- PASS: TestMustConvertInPassthrough/added_advisor (0.00s)
=== RUN   TestMustConvertInPassthrough/added_nil_item
--- PASS: TestMustConvertInPassthrough/added_nil_item (0.00s)
=== RUN   TestMustConvertInPassthrough/done_advisor
--- PASS: TestMustConvertInPassthrough/done_advisor (0.00s)
=== RUN   TestMustConvertInPassthrough/done_web_search
--- PASS: TestMustConvertInPassthrough/done_web_search (0.00s)
=== RUN   TestMustConvertInPassthrough/done_web_fetch
--- PASS: TestMustConvertInPassthrough/done_web_fetch (0.00s)
=== RUN   TestMustConvertInPassthrough/done_code_interpreter
--- PASS: TestMustConvertInPassthrough/done_code_interpreter (0.00s)
=== RUN   TestMustConvertInPassthrough/done_computer
--- PASS: TestMustConvertInPassthrough/done_computer (0.00s)
=== RUN   TestMustConvertInPassthrough/done_message
--- PASS: TestMustConvertInPassthrough/done_message (0.00s)
=== RUN   TestMustConvertInPassthrough/done_function_call
--- PASS: TestMustConvertInPassthrough/done_function_call (0.00s)
=== RUN   TestMustConvertInPassthrough/done_mcp_call
--- PASS: TestMustConvertInPassthrough/done_mcp_call (0.00s)
=== RUN   TestMustConvertInPassthrough/done_nil_item
--- PASS: TestMustConvertInPassthrough/done_nil_item (0.00s)
=== RUN   TestMustConvertInPassthrough/web_search_in_progress
--- PASS: TestMustConvertInPassthrough/web_search_in_progress (0.00s)
=== RUN   TestMustConvertInPassthrough/web_search_completed
--- PASS: TestMustConvertInPassthrough/web_search_completed (0.00s)
=== RUN   TestMustConvertInPassthrough/web_fetch_completed
--- PASS: TestMustConvertInPassthrough/web_fetch_completed (0.00s)
=== RUN   TestMustConvertInPassthrough/code_interpreter_code_done
--- PASS: TestMustConvertInPassthrough/code_interpreter_code_done (0.00s)
=== RUN   TestMustConvertInPassthrough/code_interpreter_completed
--- PASS: TestMustConvertInPassthrough/code_interpreter_completed (0.00s)
=== RUN   TestMustConvertInPassthrough/text_delta
--- PASS: TestMustConvertInPassthrough/text_delta (0.00s)
=== RUN   TestMustConvertInPassthrough/function_args_delta
--- PASS: TestMustConvertInPassthrough/function_args_delta (0.00s)
=== RUN   TestMustConvertInPassthrough/content_part_added
--- PASS: TestMustConvertInPassthrough/content_part_added (0.00s)
=== RUN   TestMustConvertInPassthrough/created
--- PASS: TestMustConvertInPassthrough/created (0.00s)
=== RUN   TestMustConvertInPassthrough/completed
--- PASS: TestMustConvertInPassthrough/completed (0.00s)
--- PASS: TestMustConvertInPassthrough (0.00s)
=== RUN   TestAnthropicContainerUploadSurvivesNormalization
--- PASS: TestAnthropicContainerUploadSurvivesNormalization (0.00s)
=== RUN   TestCheckAnthropicPassthrough_OutputConfigEscapeHatch
=== RUN   TestCheckAnthropicPassthrough_OutputConfigEscapeHatch/vertex
--- PASS: TestCheckAnthropicPassthrough_OutputConfigEscapeHatch/vertex (0.00s)
=== RUN   TestCheckAnthropicPassthrough_OutputConfigEscapeHatch/bedrock_mantle
--- PASS: TestCheckAnthropicPassthrough_OutputConfigEscapeHatch/bedrock_mantle (0.00s)
=== RUN   TestCheckAnthropicPassthrough_OutputConfigEscapeHatch/azure
--- PASS: TestCheckAnthropicPassthrough_OutputConfigEscapeHatch/azure (0.00s)
=== RUN   TestCheckAnthropicPassthrough_OutputConfigEscapeHatch/anthropic
--- PASS: TestCheckAnthropicPassthrough_OutputConfigEscapeHatch/anthropic (0.00s)
--- PASS: TestCheckAnthropicPassthrough_OutputConfigEscapeHatch (0.00s)
=== RUN   TestCheckAnthropicPassthrough_VertexBetaKeepsNativeResponseCodec
=== RUN   TestCheckAnthropicPassthrough_VertexBetaKeepsNativeResponseCodec/prompt-caching-scope-2026-01-05
--- PASS: TestCheckAnthropicPassthrough_VertexBetaKeepsNativeResponseCodec/prompt-caching-scope-2026-01-05 (0.00s)
=== RUN   TestCheckAnthropicPassthrough_VertexBetaKeepsNativeResponseCodec/fast-mode-2026-02-01
--- PASS: TestCheckAnthropicPassthrough_VertexBetaKeepsNativeResponseCodec/fast-mode-2026-02-01 (0.00s)
--- PASS: TestCheckAnthropicPassthrough_VertexBetaKeepsNativeResponseCodec (0.00s)
=== RUN   TestCheckAnthropicPassthroughLegacyCompleteOmitsStreamCodec
--- PASS: TestCheckAnthropicPassthroughLegacyCompleteOmitsStreamCodec (0.00s)
=== RUN   TestCheckAnthropicPassthrough_OAuthHeaderRouting
--- PASS: TestCheckAnthropicPassthrough_OAuthHeaderRouting (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/mantle_claude-named_alias_to_openai_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/mantle_claude-named_alias_to_openai_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/vertex_claude-named_alias_to_openai_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/vertex_claude-named_alias_to_openai_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/azure_claude-named_alias_to_openai_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/azure_claude-named_alias_to_openai_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/mantle_neutral_alias_to_claude_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/mantle_neutral_alias_to_claude_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/mantle_unaliased_claude_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/mantle_unaliased_claude_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/mantle_unaliased_openai_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/mantle_unaliased_openai_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/azure_deployment_alias_naming_claude_only_in_model_name
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/azure_deployment_alias_naming_claude_only_in_model_name (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/anthropic_provider_always_native
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/anthropic_provider_always_native (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/bedrock_is_never_native
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/bedrock_is_never_native (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_claude_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_claude_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_non-claude_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_non-claude_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_mantle_claude_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_mantle_claude_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_vertex_claude_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_vertex_claude_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_azure_claude_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_azure_claude_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_mantle_openai_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_mantle_openai_model (0.00s)
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily (0.00s)
=== RUN   TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel
=== RUN   TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel/vertex_claude
--- PASS: TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel/vertex_claude (0.00s)
=== RUN   TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel/bedrock_mantle_claude
--- PASS: TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel/bedrock_mantle_claude (0.00s)
=== RUN   TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel/azure_claude
--- PASS: TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel/azure_claude (0.00s)
=== RUN   TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel/bedrock_mantle_openai
--- PASS: TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel/bedrock_mantle_openai (0.00s)
--- PASS: TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel (0.00s)
=== RUN   TestAnthropicRawArgumentDelta
--- PASS: TestAnthropicRawArgumentDelta (0.00s)
=== RUN   TestAnthropicMessagesRawResponseCarriesExtraFields
=== RUN   TestAnthropicMessagesRawResponseCarriesExtraFields/raw_capture_requested
=== RUN   TestApplyHTTPRequestToCtxKeepsCookieOnce
--- PASS: TestApplyHTTPRequestToCtxKeepsCookieOnce (0.00s)
=== RUN   TestClearCache_OK
--- PASS: TestAnthropicMessagesRawResponseCarriesExtraFields/raw_capture_requested (0.01s)
=== RUN   TestAnthropicMessagesRawResponseCarriesExtraFields/claude_code_passthrough_stays_byte-identical
--- PASS: TestAnthropicMessagesRawResponseCarriesExtraFields/claude_code_passthrough_stays_byte-identical (0.00s)
=== RUN   TestAnthropicMessagesRawResponseCarriesExtraFields/no_raw_response_uses_the_converted_shape
--- PASS: TestAnthropicMessagesRawResponseCarriesExtraFields/no_raw_response_uses_the_converted_shape (0.00s)
--- PASS: TestAnthropicMessagesRawResponseCarriesExtraFields (0.01s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases
--- PASS: TestClearCache_OK (0.00s)
=== RUN   TestClearCache_RejectsEmptyID
--- PASS: TestClearCache_RejectsEmptyID (0.00s)
=== RUN   TestClearCache_MissingUserValue
--- PASS: TestClearCache_MissingUserValue (0.00s)
=== RUN   TestClearCache_PluginErrorReturns500
--- PASS: TestClearCache_PluginErrorReturns500 (0.00s)
=== RUN   TestClearCache_PluginNotLoaded
--- PASS: TestClearCache_PluginNotLoaded (0.00s)
=== RUN   TestClearCacheByKey_OK
--- PASS: TestClearCacheByKey_OK (0.00s)
=== RUN   TestClearCacheByKey_PluginErrorReturns500
--- PASS: TestClearCacheByKey_PluginErrorReturns500 (0.00s)
=== RUN   TestClearCacheByKey_PluginNotLoaded
--- PASS: TestClearCacheByKey_PluginNotLoaded (0.00s)
=== RUN   TestCodexWSIdentityPreservesBusinessBytes
--- PASS: TestCodexWSIdentityPreservesBusinessBytes (0.00s)
=== RUN   TestCodexWSIdentityRejectsInvalidContext
--- PASS: TestCodexWSIdentityRejectsInvalidContext (0.00s)
=== RUN   TestValidateHeaderFilterConfig
=== RUN   TestValidateHeaderFilterConfig/nil_config
--- PASS: TestValidateHeaderFilterConfig/nil_config (0.00s)
=== RUN   TestValidateHeaderFilterConfig/empty_lists
--- PASS: TestValidateHeaderFilterConfig/empty_lists (0.00s)
=== RUN   TestValidateHeaderFilterConfig/empty_allowlist_and_denylist_slices
--- PASS: TestValidateHeaderFilterConfig/empty_allowlist_and_denylist_slices (0.00s)
=== RUN   TestValidateHeaderFilterConfig/valid_allowlist_patterns
--- PASS: TestValidateHeaderFilterConfig/valid_allowlist_patterns (0.00s)
=== RUN   TestValidateHeaderFilterConfig/valid_denylist_patterns
--- PASS: TestValidateHeaderFilterConfig/valid_denylist_patterns (0.00s)
=== RUN   TestValidateHeaderFilterConfig/valid_allowlist_and_denylist_together
--- PASS: TestValidateHeaderFilterConfig/valid_allowlist_and_denylist_together (0.00s)
=== RUN   TestValidateHeaderFilterConfig/whitespace-only_entries_in_allowlist_are_dropped
--- PASS: TestValidateHeaderFilterConfig/whitespace-only_entries_in_allowlist_are_dropped (0.00s)
=== RUN   TestValidateHeaderFilterConfig/whitespace-only_entries_in_denylist_are_dropped
--- PASS: TestValidateHeaderFilterConfig/whitespace-only_entries_in_denylist_are_dropped (0.00s)
=== RUN   TestValidateHeaderFilterConfig/all-empty_allowlist_becomes_effectively_empty
--- PASS: TestValidateHeaderFilterConfig/all-empty_allowlist_becomes_effectively_empty (0.00s)
=== RUN   TestValidateHeaderFilterConfig/security_header_in_allowlist_rejected
--- PASS: TestValidateHeaderFilterConfig/security_header_in_allowlist_rejected (0.00s)
=== RUN   TestValidateHeaderFilterConfig/security_header_in_denylist_rejected
--- PASS: TestValidateHeaderFilterConfig/security_header_in_denylist_rejected (0.00s)
=== RUN   TestValidateHeaderFilterConfig/wildcard_matching_security_header_allowed_(runtime_strips_security_headers)
--- PASS: TestValidateHeaderFilterConfig/wildcard_matching_security_header_allowed_(runtime_strips_security_headers) (0.00s)
=== RUN   TestValidateHeaderFilterConfig/wildcard_prefix_matching_security_headers_allowed_(runtime_strips_security_headers)
--- PASS: TestValidateHeaderFilterConfig/wildcard_prefix_matching_security_headers_allowed_(runtime_strips_security_headers) (0.00s)
=== RUN   TestValidateHeaderFilterConfig/bare_wildcard_in_allowlist_allowed_(runtime_strips_security_headers)
--- PASS: TestValidateHeaderFilterConfig/bare_wildcard_in_allowlist_allowed_(runtime_strips_security_headers) (0.00s)
=== RUN   TestValidateHeaderFilterConfig/wildcard_in_middle_of_pattern_rejected
--- PASS: TestValidateHeaderFilterConfig/wildcard_in_middle_of_pattern_rejected (0.00s)
--- PASS: TestValidateHeaderFilterConfig (0.00s)
=== RUN   TestValidateHeaderFilterConfig_EmptyEntriesDropped
--- PASS: TestValidateHeaderFilterConfig_EmptyEntriesDropped (0.00s)
=== RUN   TestValidateHeaderFilterConfig_EmptyConfigStillForwardsHeaders
--- PASS: TestValidateHeaderFilterConfig_EmptyConfigStillForwardsHeaders (0.00s)
=== RUN   TestGetPasswordPolicyFailures
=== RUN   TestGetPasswordPolicyFailures/valid_password
--- PASS: TestGetPasswordPolicyFailures/valid_password (0.00s)
=== RUN   TestGetPasswordPolicyFailures/missing_all_requirements
--- PASS: TestGetPasswordPolicyFailures/missing_all_requirements (0.00s)
=== RUN   TestGetPasswordPolicyFailures/missing_character_classes
--- PASS: TestGetPasswordPolicyFailures/missing_character_classes (0.00s)
--- PASS: TestGetPasswordPolicyFailures (0.00s)
=== RUN   TestValidateGlobalToolSyncIntervalMinutes
--- PASS: TestValidateGlobalToolSyncIntervalMinutes (0.00s)
=== RUN   TestUpdateConfig_PersistsVKRotationCooldown
    configvkrotationcooldown_test.go:47: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/configvkrotationcooldown_test.go:47
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestUpdateConfig_PersistsVKRotationCooldown
--- FAIL: TestUpdateConfig_PersistsVKRotationCooldown (0.00s)
=== RUN   TestVirtualKeyBudgetOverrideLifecycle
    governance_test.go:235: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/governance_test.go:235
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestVirtualKeyBudgetOverrideLifecycle
--- FAIL: TestVirtualKeyBudgetOverrideLifecycle (0.00s)
=== RUN   TestVirtualKeyBudgetOverrideRejectsDirectMirrorBudget
    governance_test.go:313: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/governance_test.go:313
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestVirtualKeyBudgetOverrideRejectsDirectMirrorBudget
--- FAIL: TestVirtualKeyBudgetOverrideRejectsDirectMirrorBudget (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdatePreservesOmittedAssociation
--- PASS: TestApplyVirtualKeyOwnershipUpdatePreservesOmittedAssociation (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation
=== RUN   TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/set_team_clears_customer
--- PASS: TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/set_team_clears_customer (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/set_customer_clears_team
--- PASS: TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/set_customer_clears_team (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/null_team_clears_both
--- PASS: TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/null_team_clears_both (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/empty_customer_clears_both
--- PASS: TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/empty_customer_clears_both (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/null_team_and_customer_clears_both
--- PASS: TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/null_team_and_customer_clears_both (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/empty_team_and_customer_clears_both
--- PASS: TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/empty_team_and_customer_clears_both (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/team_with_null_customer_sets_team
--- PASS: TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/team_with_null_customer_sets_team (0.00s)
--- PASS: TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdateRejectsDualAssociation
--- PASS: TestApplyVirtualKeyOwnershipUpdateRejectsDualAssociation (0.00s)
=== RUN   TestFindExistingBudgetPrefersIDOverResetDuration
--- PASS: TestFindExistingBudgetPrefersIDOverResetDuration (0.00s)
=== RUN   TestFindExistingBudgetRejectsUnknownID
--- PASS: TestFindExistingBudgetRejectsUnknownID (0.00s)
=== RUN   TestBudgetFrequencyReplaceInheritsUsageFromOriginalBudgets
--- PASS: TestBudgetFrequencyReplaceInheritsUsageFromOriginalBudgets (0.00s)
=== RUN   TestBudgetLookupConsumesMatchedRowsForDurationSwap
--- PASS: TestBudgetLookupConsumesMatchedRowsForDurationSwap (0.00s)
=== RUN   TestResetBudgetUsageIfRequested
--- PASS: TestResetBudgetUsageIfRequested (0.00s)
=== RUN   TestTeamBudgetFrequencyChangePreservesUsageWhenRequested
--- PASS: TestTeamBudgetFrequencyChangePreservesUsageWhenRequested (0.00s)
=== RUN   TestVirtualKeyBudgetFrequencyChangePreservesUsageWhenRequested
--- PASS: TestVirtualKeyBudgetFrequencyChangePreservesUsageWhenRequested (0.00s)
=== RUN   TestProviderBudgetFrequencyChangePreservesUsageWhenRequested
--- PASS: TestProviderBudgetFrequencyChangePreservesUsageWhenRequested (0.00s)
=== RUN   TestBudgetFrequencyChangeResetsUsageWhenRequested
--- PASS: TestBudgetFrequencyChangeResetsUsageWhenRequested (0.00s)
=== RUN   TestExistingVirtualKeyBudgetLoweredBelowPreservedUsageIsAllowed
--- PASS: TestExistingVirtualKeyBudgetLoweredBelowPreservedUsageIsAllowed (0.00s)
=== RUN   TestExistingProviderBudgetLoweredBelowPreservedUsageIsAllowed
--- PASS: TestExistingProviderBudgetLoweredBelowPreservedUsageIsAllowed (0.00s)
=== RUN   TestExistingBudgetLoweredBelowUsageSucceedsWhenResetRequested
--- PASS: TestExistingBudgetLoweredBelowUsageSucceedsWhenResetRequested (0.00s)
=== RUN   TestNewVirtualKeyBudgetInheritsClosestShorterUsage
--- PASS: TestNewVirtualKeyBudgetInheritsClosestShorterUsage (0.00s)
=== RUN   TestNewProviderBudgetInheritsClosestShorterUsage
--- PASS: TestNewProviderBudgetInheritsClosestShorterUsage (0.00s)
=== RUN   TestNewVirtualKeyBudgetInheritanceAboveLimitIsAllowed
--- PASS: TestNewVirtualKeyBudgetInheritanceAboveLimitIsAllowed (0.00s)
=== RUN   TestNewProviderBudgetInheritanceAtLimitIsAllowed
--- PASS: TestNewProviderBudgetInheritanceAtLimitIsAllowed (0.00s)
=== RUN   TestNewShorterBudgetDoesNotInheritFromLongerUsage
--- PASS: TestNewShorterBudgetDoesNotInheritFromLongerUsage (0.00s)
=== RUN   TestNewBudgetInheritsClosestShorterUsage
--- PASS: TestNewBudgetInheritsClosestShorterUsage (0.00s)
=== RUN   TestNewLongerBudgetDoesNotInheritUsageWhenResetRequested
--- PASS: TestNewLongerBudgetDoesNotInheritUsageWhenResetRequested (0.00s)
=== RUN   TestRotateVirtualKey_OnlyChangesValueAndReloads
--- PASS: TestRotateVirtualKey_OnlyChangesValueAndReloads (0.00s)
=== RUN   TestRotateVirtualKey_CooldownStoresPreviousValue
--- PASS: TestRotateVirtualKey_CooldownStoresPreviousValue (0.00s)
=== RUN   TestRotateVirtualKey_ZeroCooldownClearsPreviousValue
--- PASS: TestRotateVirtualKey_ZeroCooldownClearsPreviousValue (0.00s)
=== RUN   TestRotateVirtualKey_DefaultUnsetCooldownRevokesImmediately
--- PASS: TestRotateVirtualKey_DefaultUnsetCooldownRevokesImmediately (0.00s)
=== RUN   TestRotateVirtualKey_NoClientConfigRevokesImmediately
--- PASS: TestRotateVirtualKey_NoClientConfigRevokesImmediately (0.00s)
=== RUN   TestRotateVirtualKey_ClientConfigErrorRevokesImmediately
--- PASS: TestRotateVirtualKey_ClientConfigErrorRevokesImmediately (0.00s)
=== RUN   TestRotateVirtualKeys_BulkDefaultCooldownRevokesImmediately
--- PASS: TestRotateVirtualKeys_BulkDefaultCooldownRevokesImmediately (0.00s)
=== RUN   TestRotateVirtualKey_NotFound
--- PASS: TestRotateVirtualKey_NotFound (0.00s)
=== RUN   TestRotateVirtualKey_UpdateFailureDoesNotReload
--- PASS: TestRotateVirtualKey_UpdateFailureDoesNotReload (0.00s)
=== RUN   TestRotateVirtualKey_ReloadFailureReturnsErrorAfterUpdate
--- PASS: TestRotateVirtualKey_ReloadFailureReturnsErrorAfterUpdate (0.00s)
=== RUN   TestRotateVirtualKeys_PartialSuccess
--- PASS: TestRotateVirtualKeys_PartialSuccess (0.00s)
=== RUN   TestRotateVirtualKeys_RejectsInvalidRequests
=== RUN   TestRotateVirtualKeys_RejectsInvalidRequests/invalid_JSON
--- PASS: TestRotateVirtualKeys_RejectsInvalidRequests/invalid_JSON (0.00s)
=== RUN   TestRotateVirtualKeys_RejectsInvalidRequests/empty_IDs
--- PASS: TestRotateVirtualKeys_RejectsInvalidRequests/empty_IDs (0.00s)
=== RUN   TestRotateVirtualKeys_RejectsInvalidRequests/blank_ID
--- PASS: TestRotateVirtualKeys_RejectsInvalidRequests/blank_ID (0.00s)
--- PASS: TestRotateVirtualKeys_RejectsInvalidRequests (0.00s)
=== RUN   TestRotateVirtualKeys_TrimsAndDeduplicatesIDs
--- PASS: TestRotateVirtualKeys_TrimsAndDeduplicatesIDs (0.00s)
=== RUN   TestRotateVirtualKeys_AllFailuresReturnsServerError
--- PASS: TestRotateVirtualKeys_AllFailuresReturnsServerError (0.00s)
=== RUN   TestGetVirtualKeyQuota_HydratesBudgetsFromModelConfigs
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/continue_is_refused
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/continue_is_refused (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/create_falls_through
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/create_falls_through (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/no_thread_falls_through
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/no_thread_falls_through (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/malformed_thread_falls_through
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/malformed_thread_falls_through (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/nil_extra_params_falls_through
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/nil_extra_params_falls_through (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/large_payload_continue_refused_from_metadata
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/large_payload_continue_refused_from_metadata (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/large_payload_create_falls_through
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/large_payload_create_falls_through (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/large_payload_without_thread_metadata_falls_through
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/large_payload_without_thread_metadata_falls_through (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/parsed_thread_wins_over_metadata
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/parsed_thread_wins_over_metadata (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/wrong_request_type_falls_through
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/wrong_request_type_falls_through (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/count_tokens_path_never_refused
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/count_tokens_path_never_refused (0.00s)
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases (0.02s)
=== RUN   TestGenericRouter_MarkDeprecatedListModelsResponseUsesCatalog
--- PASS: TestGetVirtualKeyQuota_HydratesBudgetsFromModelConfigs (0.01s)
=== RUN   TestGetVirtualKeyQuota_ExternalResolverReplacesWithAccessProfileBudgets
--- PASS: TestGetVirtualKeyQuota_ExternalResolverReplacesWithAccessProfileBudgets (0.00s)
=== RUN   TestGetVirtualKeyQuota_ExternalResolverRateLimitOnly
--- PASS: TestGetVirtualKeyQuota_ExternalResolverRateLimitOnly (0.00s)
=== RUN   TestApplyExternalBudgets_RateLimitOnlyDropsNativeBudgets
--- PASS: TestApplyExternalBudgets_RateLimitOnlyDropsNativeBudgets (0.00s)
=== RUN   TestApplyExternalBudgets_ManagedWithNoGovernanceFlagsAndClears
--- PASS: TestApplyExternalBudgets_ManagedWithNoGovernanceFlagsAndClears (0.00s)
=== RUN   TestGetVirtualKeyQuota_ExternalResolverErrorFailsClosed
--- PASS: TestGetVirtualKeyQuota_ExternalResolverErrorFailsClosed (0.00s)
=== RUN   TestGetVirtualKeyQuota_NoGovernanceReturnsEmpty
--- PASS: TestGetVirtualKeyQuota_NoGovernanceReturnsEmpty (0.00s)
=== RUN   TestGetVirtualKeyQuota_MissingHeaderReturns401
--- PASS: TestGetVirtualKeyQuota_MissingHeaderReturns401 (0.00s)
=== RUN   TestGetVirtualKeyQuota_NotFoundReturns401
--- PASS: TestGetVirtualKeyQuota_NotFoundReturns401 (0.00s)
=== RUN   TestGetVirtualKeyQuota_ModelConfigLoadErrorFailsClosed
--- PASS: TestGetVirtualKeyQuota_ModelConfigLoadErrorFailsClosed (0.00s)
=== RUN   TestGetVirtualKeyQuota_RankingsErrorFailsClosed
--- PASS: TestGetVirtualKeyQuota_RankingsErrorFailsClosed (0.00s)
=== RUN   TestGetVirtualKeyQuota_EndToEndWithRealStore
    governance_test.go:2309: failed to create config store: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
--- FAIL: TestGetVirtualKeyQuota_EndToEndWithRealStore (0.00s)
=== RUN   TestGetVirtualKeyQuota_GraceValueWithRealStore
    governance_test.go:2499: failed to create config store: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
--- FAIL: TestGetVirtualKeyQuota_GraceValueWithRealStore (0.00s)
=== RUN   TestGetVirtualKeyQuota_ExpiredGraceValueUnauthorized
    governance_test.go:2524: failed to create config store: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
--- FAIL: TestGetVirtualKeyQuota_ExpiredGraceValueUnauthorized (0.00s)
=== RUN   TestGetVirtualKeyQuota_ExpiredVirtualKeyRejected
--- PASS: TestGetVirtualKeyQuota_ExpiredVirtualKeyRejected (0.00s)
=== RUN   TestGetVirtualKeyQuota_UnexpiredVirtualKeyAllowed
--- PASS: TestGetVirtualKeyQuota_UnexpiredVirtualKeyAllowed (0.00s)
=== RUN   TestGetVirtualKeyQuota_InactiveVirtualKeyStillReadsQuota
--- PASS: TestGetVirtualKeyQuota_InactiveVirtualKeyStillReadsQuota (0.00s)
=== RUN   TestGetVirtualKeyQuota_WindowClampedToBudgetCreation
    governance_test.go:2646: failed to create config store: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
--- FAIL: TestGetVirtualKeyQuota_WindowClampedToBudgetCreation (0.00s)
=== RUN   TestGetVirtualKeys_PaginatedEndpoint_ResponseShape
--- PASS: TestGetVirtualKeys_PaginatedEndpoint_ResponseShape (0.00s)
=== RUN   TestGetVirtualKeys_PaginatedEndpoint_QueryParams
=== RUN   TestGetVirtualKeys_PaginatedEndpoint_QueryParams/explicit_limit_and_offset
--- PASS: TestGetVirtualKeys_PaginatedEndpoint_QueryParams/explicit_limit_and_offset (0.00s)
=== RUN   TestGetVirtualKeys_PaginatedEndpoint_QueryParams/no_params_uses_defaults
--- PASS: TestGetVirtualKeys_PaginatedEndpoint_QueryParams/no_params_uses_defaults (0.00s)
--- PASS: TestGetVirtualKeys_PaginatedEndpoint_QueryParams (0.00s)
=== RUN   TestGetVirtualKeys_FromMemoryUsesGovernanceData
--- PASS: TestGetVirtualKeys_FromMemoryUsesGovernanceData (0.00s)
=== RUN   TestGetVirtualKeys_FromMemoryTakesPrecedenceOverLimit
--- PASS: TestGetVirtualKeys_FromMemoryTakesPrecedenceOverLimit (0.00s)
=== RUN   TestGetVirtualKeys_FromMemoryRejectsUserFilter
--- PASS: TestGetVirtualKeys_FromMemoryRejectsUserFilter (0.00s)
=== RUN   TestGetVirtualKeys_FromMemoryIgnoresOtherFilters
--- PASS: TestGetVirtualKeys_FromMemoryIgnoresOtherFilters (0.00s)
=== RUN   TestBudgetRemovalRequestDetection
=== RUN   TestBudgetRemovalRequestDetection/nil_request_is_not_removal
--- PASS: TestBudgetRemovalRequestDetection/nil_request_is_not_removal (0.00s)
=== RUN   TestBudgetRemovalRequestDetection/empty_object_is_removal
--- PASS: TestBudgetRemovalRequestDetection/empty_object_is_removal (0.00s)
=== RUN   TestBudgetRemovalRequestDetection/max_limit_present_is_not_removal
--- PASS: TestBudgetRemovalRequestDetection/max_limit_present_is_not_removal (0.00s)
=== RUN   TestBudgetRemovalRequestDetection/reset_duration_only_is_not_removal
--- PASS: TestBudgetRemovalRequestDetection/reset_duration_only_is_not_removal (0.00s)
--- PASS: TestBudgetRemovalRequestDetection (0.00s)
=== RUN   TestRateLimitRemovalRequestDetection
=== RUN   TestRateLimitRemovalRequestDetection/nil_request_is_not_removal
--- PASS: TestRateLimitRemovalRequestDetection/nil_request_is_not_removal (0.00s)
=== RUN   TestRateLimitRemovalRequestDetection/empty_object_is_removal
--- PASS: TestRateLimitRemovalRequestDetection/empty_object_is_removal (0.00s)
=== RUN   TestRateLimitRemovalRequestDetection/token_limit_present_is_not_removal
--- PASS: TestRateLimitRemovalRequestDetection/token_limit_present_is_not_removal (0.00s)
=== RUN   TestRateLimitRemovalRequestDetection/request_limit_present_is_not_removal
--- PASS: TestRateLimitRemovalRequestDetection/request_limit_present_is_not_removal (0.00s)
=== RUN   TestRateLimitRemovalRequestDetection/durations_only_is_not_removal
--- PASS: TestRateLimitRemovalRequestDetection/durations_only_is_not_removal (0.00s)
--- PASS: TestRateLimitRemovalRequestDetection (0.00s)
=== RUN   TestCollectProviderConfigDeleteIDs
=== RUN   TestCollectProviderConfigDeleteIDs/collects_both_IDs
--- PASS: TestCollectProviderConfigDeleteIDs/collects_both_IDs (0.00s)
=== RUN   TestCollectProviderConfigDeleteIDs/appends_to_existing_slices
--- PASS: TestCollectProviderConfigDeleteIDs/appends_to_existing_slices (0.00s)
=== RUN   TestCollectProviderConfigDeleteIDs/ignores_missing_IDs
--- PASS: TestCollectProviderConfigDeleteIDs/ignores_missing_IDs (0.00s)
--- PASS: TestCollectProviderConfigDeleteIDs (0.00s)
=== RUN   TestCoerceLegacyBudget
=== RUN   TestCoerceLegacyBudget/empty_object_→_removal,_returns_empty_slice
--- PASS: TestCoerceLegacyBudget/empty_object_→_removal,_returns_empty_slice (0.00s)
=== RUN   TestCoerceLegacyBudget/both_fields_set,_no_existing_→_new_budget_entry,_no_ID
--- PASS: TestCoerceLegacyBudget/both_fields_set,_no_existing_→_new_budget_entry,_no_ID (0.00s)
=== RUN   TestCoerceLegacyBudget/update_max_limit_only,_existing_budget_→_merges_ID_and_reset_duration
--- PASS: TestCoerceLegacyBudget/update_max_limit_only,_existing_budget_→_merges_ID_and_reset_duration (0.00s)
=== RUN   TestCoerceLegacyBudget/update_reset_duration_only,_existing_budget_→_merges_ID_and_max_limit
--- PASS: TestCoerceLegacyBudget/update_reset_duration_only,_existing_budget_→_merges_ID_and_max_limit (0.00s)
=== RUN   TestCoerceLegacyBudget/max_limit_only,_no_existing_→_cannot_build_valid_budget,_returns_nil
--- PASS: TestCoerceLegacyBudget/max_limit_only,_no_existing_→_cannot_build_valid_budget,_returns_nil (0.00s)
=== RUN   TestCoerceLegacyBudget/reset_duration_only,_no_existing_→_cannot_build_valid_budget,_returns_nil
--- PASS: TestCoerceLegacyBudget/reset_duration_only,_no_existing_→_cannot_build_valid_budget,_returns_nil (0.00s)
--- PASS: TestCoerceLegacyBudget (0.00s)
=== RUN   TestModelConfigToProviderGovernanceNewFields
=== RUN   TestModelConfigToProviderGovernanceNewFields/nil_mc_returns_false
--- PASS: TestModelConfigToProviderGovernanceNewFields/nil_mc_returns_false (0.00s)
=== RUN   TestModelConfigToProviderGovernanceNewFields/wrong_scope_returns_false
--- PASS: TestModelConfigToProviderGovernanceNewFields/wrong_scope_returns_false (0.00s)
=== RUN   TestModelConfigToProviderGovernanceNewFields/no_budgets:_Budget_nil,_Budgets_empty,_CalendarAligned_false
--- PASS: TestModelConfigToProviderGovernanceNewFields/no_budgets:_Budget_nil,_Budgets_empty,_CalendarAligned_false (0.00s)
=== RUN   TestModelConfigToProviderGovernanceNewFields/single_budget:_Budget_points_to_first,_Budgets_has_one_entry
--- PASS: TestModelConfigToProviderGovernanceNewFields/single_budget:_Budget_points_to_first,_Budgets_has_one_entry (0.00s)
=== RUN   TestModelConfigToProviderGovernanceNewFields/multiple_budgets:_Budget_is_first,_Budgets_contains_all
--- PASS: TestModelConfigToProviderGovernanceNewFields/multiple_budgets:_Budget_is_first,_Budgets_contains_all (0.00s)
=== RUN   TestModelConfigToProviderGovernanceNewFields/calendar_aligned_is_propagated
--- PASS: TestModelConfigToProviderGovernanceNewFields/calendar_aligned_is_propagated (0.00s)
=== RUN   TestModelConfigToProviderGovernanceNewFields/Budgets_slice_is_a_copy,_not_a_reference_to_mc.Budgets
--- PASS: TestModelConfigToProviderGovernanceNewFields/Budgets_slice_is_a_copy,_not_a_reference_to_mc.Budgets (0.00s)
--- PASS: TestModelConfigToProviderGovernanceNewFields (0.00s)
=== RUN   TestUpdateProviderGovernance_BudgetMutualExclusion
--- PASS: TestUpdateProviderGovernance_BudgetMutualExclusion (0.00s)
=== RUN   TestValidateRoutingFallbacks
=== RUN   TestValidateRoutingFallbacks/nil
--- PASS: TestValidateRoutingFallbacks/nil (0.00s)
=== RUN   TestValidateRoutingFallbacks/empty
--- PASS: TestValidateRoutingFallbacks/empty (0.00s)
=== RUN   TestValidateRoutingFallbacks/provider_model
--- PASS: TestValidateRoutingFallbacks/provider_model (0.00s)
=== RUN   TestValidateRoutingFallbacks/provider_slash_incoming_model
--- PASS: TestValidateRoutingFallbacks/provider_slash_incoming_model (0.00s)
=== RUN   TestValidateRoutingFallbacks/bare_known_provider_name_rejected
--- PASS: TestValidateRoutingFallbacks/bare_known_provider_name_rejected (0.00s)
=== RUN   TestValidateRoutingFallbacks/bare_model_rejected
--- PASS: TestValidateRoutingFallbacks/bare_model_rejected (0.00s)
=== RUN   TestValidateRoutingFallbacks/empty_element
--- PASS: TestValidateRoutingFallbacks/empty_element (0.00s)
=== RUN   TestValidateRoutingFallbacks/huggingface_namespace_not_a_provider_prefix
--- PASS: TestValidateRoutingFallbacks/huggingface_namespace_not_a_provider_prefix (0.00s)
--- PASS: TestValidateRoutingFallbacks (0.00s)
=== RUN   TestCreateCustomer_CalendarAligned_SnapsBudgetLastReset
--- PASS: TestCreateCustomer_CalendarAligned_SnapsBudgetLastReset (0.00s)
=== RUN   TestCreateCustomer_CalendarAligned_False
--- PASS: TestCreateCustomer_CalendarAligned_False (0.00s)
=== RUN   TestUpdateCustomer_CalendarAligned_DoesNotTouchBudgets
--- PASS: TestUpdateCustomer_CalendarAligned_DoesNotTouchBudgets (0.01s)
=== RUN   TestUpdateCustomer_CalendarAligned_NoSnapWhenAlreadyEnabled
--- PASS: TestUpdateCustomer_CalendarAligned_NoSnapWhenAlreadyEnabled (0.00s)
=== RUN   TestApplyVKGovernanceFromModelConfigs_PreservesDirectlyAttachedBudget
--- PASS: TestApplyVKGovernanceFromModelConfigs_PreservesDirectlyAttachedBudget (0.00s)
=== RUN   TestApplyVKGovernanceFromModelConfigs_OverlaysModelConfigGovernance
--- PASS: TestApplyVKGovernanceFromModelConfigs_OverlaysModelConfigGovernance (0.00s)
=== RUN   TestProviderGovernance_DecodesEncodedProviderName
    governance_test.go:3752: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/governance_test.go:3752
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestProviderGovernance_DecodesEncodedProviderName
--- FAIL: TestProviderGovernance_DecodesEncodedProviderName (0.00s)
=== RUN   TestProviderGovernance_UnknownProviderStill404
    governance_test.go:3802: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/governance_test.go:3802
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestProviderGovernance_UnknownProviderStill404
--- FAIL: TestProviderGovernance_UnknownProviderStill404 (0.00s)
=== RUN   TestUpdateProviderGovernance_AdoptsReconciledBudgetsNotStaleOnes
    governance_test.go:3846: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/governance_test.go:3846
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestUpdateProviderGovernance_AdoptsReconciledBudgetsNotStaleOnes
--- FAIL: TestUpdateProviderGovernance_AdoptsReconciledBudgetsNotStaleOnes (0.00s)
=== RUN   TestProviderGovernance_MalformedEncodingReturns400
    governance_test.go:3895: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/governance_test.go:3895
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestProviderGovernance_MalformedEncodingReturns400
--- FAIL: TestProviderGovernance_MalformedEncodingReturns400 (0.00s)
=== RUN   TestTeam_DecodesEncodedTeamID
    governance_test.go:3933: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/governance_test.go:3933
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestTeam_DecodesEncodedTeamID
--- FAIL: TestTeam_DecodesEncodedTeamID (0.00s)
=== RUN   TestTeam_MalformedEncodingReturns400
    governance_test.go:4013: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/governance_test.go:4013
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestTeam_MalformedEncodingReturns400
--- FAIL: TestTeam_MalformedEncodingReturns400 (0.00s)
=== RUN   TestValidateBudgetResetConfig
=== RUN   TestValidateBudgetResetConfig/quarterly_with_no_config_is_valid
--- PASS: TestValidateBudgetResetConfig/quarterly_with_no_config_is_valid (0.00s)
=== RUN   TestValidateBudgetResetConfig/quarterly_with_a_fiscal_start_is_valid
--- PASS: TestValidateBudgetResetConfig/quarterly_with_a_fiscal_start_is_valid (0.00s)
=== RUN   TestValidateBudgetResetConfig/monthly_with_a_quarter_config_is_rejected
--- PASS: TestValidateBudgetResetConfig/monthly_with_a_quarter_config_is_rejected (0.00s)
=== RUN   TestValidateBudgetResetConfig/hourly_with_a_quarter_config_is_rejected
--- PASS: TestValidateBudgetResetConfig/hourly_with_a_quarter_config_is_rejected (0.00s)
=== RUN   TestValidateBudgetResetConfig/month_above_december_is_rejected
--- PASS: TestValidateBudgetResetConfig/month_above_december_is_rejected (0.00s)
=== RUN   TestValidateBudgetResetConfig/negative_month_is_rejected
--- PASS: TestValidateBudgetResetConfig/negative_month_is_rejected (0.00s)
=== RUN   TestValidateBudgetResetConfig/zero_month_means_unset_and_is_accepted
--- PASS: TestValidateBudgetResetConfig/zero_month_means_unset_and_is_accepted (0.00s)
--- PASS: TestValidateBudgetResetConfig (0.00s)
=== RUN   TestNewBudgetFromRequestCarriesResetConfig
--- PASS: TestNewBudgetFromRequestCarriesResetConfig (0.00s)
=== RUN   TestNewBudgetFromRequestSnapsToFiscalQuarter
--- PASS: TestNewBudgetFromRequestSnapsToFiscalQuarter (0.00s)
=== RUN   TestApplyResetConfigToExistingBudgetIsConfigOnly
--- PASS: TestApplyResetConfigToExistingBudgetIsConfigOnly (0.00s)
=== RUN   TestApplyResetConfigToExistingBudgetClearsDefinition
--- PASS: TestApplyResetConfigToExistingBudgetClearsDefinition (0.00s)
=== RUN   TestQuarterStartChangeMakesBudgetDue
--- PASS: TestQuarterStartChangeMakesBudgetDue (0.00s)
=== RUN   TestBudgetLastResetUsesBudgetQuarterStart
--- PASS: TestBudgetLastResetUsesBudgetQuarterStart (0.00s)
=== RUN   TestApplyAssignees
=== RUN   TestApplyAssignees/fills_in_assignees_in_one_batched_call
--- PASS: TestApplyAssignees/fills_in_assignees_in_one_batched_call (0.00s)
=== RUN   TestApplyAssignees/no-ops_without_a_resolver
--- PASS: TestApplyAssignees/no-ops_without_a_resolver (0.00s)
=== RUN   TestApplyAssignees/degrades_to_no_assignee_when_the_resolver_fails
--- PASS: TestApplyAssignees/degrades_to_no_assignee_when_the_resolver_fails (0.00s)
=== RUN   TestApplyAssignees/skips_the_resolver_for_an_empty_page
--- PASS: TestApplyAssignees/skips_the_resolver_for_an_empty_page (0.00s)
--- PASS: TestApplyAssignees (0.00s)
=== RUN   TestGovernanceRouteOverridesReplaceTeamFamily
--- PASS: TestGovernanceRouteOverridesReplaceTeamFamily (0.00s)
=== RUN   TestGovernanceRouteOverridesDefaultToOSS
--- PASS: TestGovernanceRouteOverridesDefaultToOSS (0.00s)
=== RUN   TestResolveBatchProvider
=== RUN   TestResolveBatchProvider/model_field:_provider+model_parsed
--- PASS: TestResolveBatchProvider/model_field:_provider+model_parsed (0.00s)
=== RUN   TestResolveBatchProvider/no_model,_x-model-provider_header
--- PASS: TestResolveBatchProvider/no_model,_x-model-provider_header (0.00s)
=== RUN   TestResolveBatchProvider/no_model,_?provider=_query_param
--- PASS: TestResolveBatchProvider/no_model,_?provider=_query_param (0.00s)
=== RUN   TestResolveBatchProvider/no_model,_no_provider_→_error
--- PASS: TestResolveBatchProvider/no_model,_no_provider_→_error (0.00s)
--- PASS: TestResolveBatchProvider (0.00s)
=== RUN   TestPrepareImageEditRequest_JSON
--- PASS: TestPrepareImageEditRequest_JSON (0.01s)
=== RUN   TestPrepareImageEditRequest_JSONTypedExtraParams
--- PASS: TestPrepareImageEditRequest_JSONTypedExtraParams (0.00s)
=== RUN   TestPrepareImageEditRequest_JSONInputFieldsAreNotExtraParams
--- PASS: TestPrepareImageEditRequest_JSONInputFieldsAreNotExtraParams (0.00s)
=== RUN   TestPrepareImageEditRequest_JSONBareStringImages
--- PASS: TestPrepareImageEditRequest_JSONBareStringImages (0.00s)
=== RUN   TestPrepareImageEditRequest_JSONBase64Image
--- PASS: TestPrepareImageEditRequest_JSONBase64Image (0.00s)
=== RUN   TestPrepareImageEditRequest_MultipartUnchanged
--- PASS: TestPrepareImageEditRequest_MultipartUnchanged (0.00s)
=== RUN   TestPrepareImageEditRequest_MultipartUpload
--- PASS: TestPrepareImageEditRequest_MultipartUpload (0.00s)
=== RUN   TestPrepareImageEditRequest_Validation
--- PASS: TestPrepareImageEditRequest_Validation (0.00s)
=== RUN   TestPrepareImageEditRequest_EmptyImageEntriesAreDropped
--- PASS: TestPrepareImageEditRequest_EmptyImageEntriesAreDropped (0.00s)
=== RUN   TestPrepareImageEditRequest_Stream
--- PASS: TestPrepareImageEditRequest_Stream (0.00s)
=== RUN   TestApplyListModelsProviderFilterDelegatesToTheModelsManager
--- PASS: TestApplyListModelsProviderFilterDelegatesToTheModelsManager (0.00s)
=== RUN   TestApplyListModelsProviderFilterLeavesFanOutAloneWhenNothingResolved
--- PASS: TestApplyListModelsProviderFilterLeavesFanOutAloneWhenNothingResolved (0.00s)
=== RUN   TestApplyListModelsProviderFilterWithoutModelsManager
--- PASS: TestApplyListModelsProviderFilterWithoutModelsManager (0.00s)
=== RUN   TestPrepareVideoEditRequest_JSON
--- PASS: TestPrepareVideoEditRequest_JSON (0.00s)
=== RUN   TestPrepareVideoEditRequest_ProviderFromVideoID
--- PASS: TestPrepareVideoEditRequest_ProviderFromVideoID (0.00s)
=== RUN   TestPrepareVideoEditRequest_ProviderFallbacks
--- PASS: TestPrepareVideoEditRequest_ProviderFallbacks (0.00s)
=== RUN   TestPrepareVideoEditRequest_BareModelStillResolvesProvider
--- PASS: TestPrepareVideoEditRequest_BareModelStillResolvesProvider (0.00s)
=== RUN   TestPrepareVideoEditRequest_BareModelDefersToRouting
--- PASS: TestPrepareVideoEditRequest_BareModelDefersToRouting (0.00s)
=== RUN   TestVideoIDProviderSuffix
--- PASS: TestVideoIDProviderSuffix (0.00s)
=== RUN   TestPrepareVideoEditRequest_RequiresSource
--- PASS: TestPrepareVideoEditRequest_RequiresSource (0.00s)
=== RUN   TestPrepareVideoEditRequest_ExtraParams
--- PASS: TestPrepareVideoEditRequest_ExtraParams (0.00s)
=== RUN   TestPrepareVideoEditRequest_MultipartUpload
--- PASS: TestPrepareVideoEditRequest_MultipartUpload (0.00s)
=== RUN   TestPrepareVideoEditRequest_MultipartBracketedVideoID
--- PASS: TestPrepareVideoEditRequest_MultipartBracketedVideoID (0.00s)
=== RUN   TestVideoRouteShapesDoNotConflict
--- PASS: TestVideoRouteShapesDoNotConflict (0.00s)
=== RUN   TestIsLocalhost
--- PASS: TestIsLocalhost (0.00s)
=== RUN   TestIsLocalhostOrigin
--- PASS: TestIsLocalhostOrigin (0.00s)
=== RUN   TestLoopbackRedirectURIsIPv6
--- PASS: TestLoopbackRedirectURIsIPv6 (0.00s)
=== RUN   TestPrivateUseRedirectSchemes
--- PASS: TestPrivateUseRedirectSchemes (0.00s)
=== RUN   TestMatchRedirectURIPrivateUseSchemes
--- PASS: TestMatchRedirectURIPrivateUseSchemes (0.00s)
=== RUN   TestShouldUseFilterDataCacheAllowsUnscopedEmptyQuery
--- PASS: TestShouldUseFilterDataCacheAllowsUnscopedEmptyQuery (0.00s)
=== RUN   TestParseComplexityFilters
=== RUN   TestParseComplexityFilters/parses_tier_and_mechanism
--- PASS: TestParseComplexityFilters/parses_tier_and_mechanism (0.00s)
=== RUN   TestParseComplexityFilters/leaves_filters_unchanged_when_parameters_are_absent
--- PASS: TestParseComplexityFilters/leaves_filters_unchanged_when_parameters_are_absent (0.00s)
--- PASS: TestParseComplexityFilters (0.00s)
=== RUN   TestParseToolCallNamesFilter
=== RUN   TestParseToolCallNamesFilter/parses_comma-separated_names
--- PASS: TestParseToolCallNamesFilter/parses_comma-separated_names (0.00s)
=== RUN   TestParseToolCallNamesFilter/leaves_filters_unchanged_when_parameter_is_absent
--- PASS: TestParseToolCallNamesFilter/leaves_filters_unchanged_when_parameter_is_absent (0.00s)
=== RUN   TestParseToolCallNamesFilter/histogram_filters_honour_it
--- PASS: TestParseToolCallNamesFilter/histogram_filters_honour_it (0.00s)
--- PASS: TestParseToolCallNamesFilter (0.00s)
=== RUN   TestParseParentRequestIDFilter
--- PASS: TestParseParentRequestIDFilter (0.00s)
=== RUN   TestShouldUseFilterDataCacheRejectsSearchQuery
--- PASS: TestShouldUseFilterDataCacheRejectsSearchQuery (0.00s)
=== RUN   TestShouldUseFilterDataCacheRejectsScopedContext
--- PASS: TestShouldUseFilterDataCacheRejectsScopedContext (0.00s)
=== RUN   TestGetMCPLogByIDRedactionMapping
=== RUN   TestGetMCPLogByIDRedactionMapping/no_resolver
--- PASS: TestGetMCPLogByIDRedactionMapping/no_resolver (0.01s)
=== RUN   TestGetMCPLogByIDRedactionMapping/authorized_mapping
--- PASS: TestGetMCPLogByIDRedactionMapping/authorized_mapping (0.00s)
=== RUN   TestGetMCPLogByIDRedactionMapping/resolver_error
--- PASS: TestGetMCPLogByIDRedactionMapping/resolver_error (0.00s)
--- PASS: TestGetMCPLogByIDRedactionMapping (0.01s)
=== RUN   TestShouldCacheFilterDimensions_NarrowsToRawScans
=== RUN   TestShouldCacheFilterDimensions_NarrowsToRawScans/single_matview_dimension
--- PASS: TestShouldCacheFilterDimensions_NarrowsToRawScans/single_matview_dimension (0.00s)
=== RUN   TestShouldCacheFilterDimensions_NarrowsToRawScans/several_matview_dimensions
--- PASS: TestShouldCacheFilterDimensions_NarrowsToRawScans/several_matview_dimensions (0.00s)
=== RUN   TestShouldCacheFilterDimensions_NarrowsToRawScans/metadata_keys_alone
--- PASS: TestShouldCacheFilterDimensions_NarrowsToRawScans/metadata_keys_alone (0.00s)
=== RUN   TestShouldCacheFilterDimensions_NarrowsToRawScans/metadata_keys_mixed_in
--- PASS: TestShouldCacheFilterDimensions_NarrowsToRawScans/metadata_keys_mixed_in (0.00s)
=== RUN   TestShouldCacheFilterDimensions_NarrowsToRawScans/default_all_dimensions
--- PASS: TestShouldCacheFilterDimensions_NarrowsToRawScans/default_all_dimensions (0.00s)
--- PASS: TestShouldCacheFilterDimensions_NarrowsToRawScans (0.00s)
=== RUN   TestShouldCacheFilterDimensions_NonPostgresCachesEverything
--- PASS: TestShouldCacheFilterDimensions_NonPostgresCachesEverything (0.00s)
=== RUN   TestFilterDataCacheIdentity_PartitionsPerCaller
--- PASS: TestFilterDataCacheIdentity_PartitionsPerCaller (0.00s)
=== RUN   TestGetDashboard
=== RUN   TestGetDashboard/success_includes_all_sections
--- PASS: TestGetDashboard/success_includes_all_sections (0.02s)
=== RUN   TestGetDashboard/sub-query_error_returns_no_partial_dashboard
--- PASS: TestGetDashboard/sub-query_error_returns_no_partial_dashboard (0.00s)
=== RUN   TestGetDashboard/MCP_filters_are_isolated_from_LLM_filters
--- PASS: TestGetDashboard/MCP_filters_are_isolated_from_LLM_filters (0.00s)
--- PASS: TestGetDashboard (0.02s)
=== RUN   TestRecalculateLogCostsResolvesPeriodFilter
--- PASS: TestRecalculateLogCostsResolvesPeriodFilter (0.03s)
=== RUN   TestRecalculateLogCostsRejectsDuplicateJob
--- PASS: TestRecalculateLogCostsRejectsDuplicateJob (0.00s)
=== RUN   TestCancelRecalculateCost
=== RUN   TestCancelRecalculateCost/cancels_the_in-flight_job_when_no_id_is_given
--- PASS: TestCancelRecalculateCost/cancels_the_in-flight_job_when_no_id_is_given (0.00s)
=== RUN   TestCancelRecalculateCost/cancels_the_job_named_by_id
--- PASS: TestCancelRecalculateCost/cancels_the_job_named_by_id (0.00s)
=== RUN   TestCancelRecalculateCost/an_already-terminal_job_is_returned_unchanged
--- PASS: TestCancelRecalculateCost/an_already-terminal_job_is_returned_unchanged (0.00s)
=== RUN   TestCancelRecalculateCost/refuses_to_cancel_a_job_of_another_kind
--- PASS: TestCancelRecalculateCost/refuses_to_cancel_a_job_of_another_kind (0.00s)
=== RUN   TestCancelRecalculateCost/404_when_there_is_nothing_to_cancel
--- PASS: TestCancelRecalculateCost/404_when_there_is_nothing_to_cancel (0.00s)
=== RUN   TestCancelRecalculateCost/503_when_the_background_runner_is_not_wired
--- PASS: TestCancelRecalculateCost/503_when_the_background_runner_is_not_wired (0.00s)
--- PASS: TestCancelRecalculateCost (0.00s)
=== RUN   TestProjectFilterReachesEveryLogQuery
=== RUN   TestProjectFilterReachesEveryLogQuery/list
--- PASS: TestProjectFilterReachesEveryLogQuery/list (0.00s)
=== RUN   TestProjectFilterReachesEveryLogQuery/stats
--- PASS: TestProjectFilterReachesEveryLogQuery/stats (0.00s)
=== RUN   TestProjectFilterReachesEveryLogQuery/histograms
--- PASS: TestProjectFilterReachesEveryLogQuery/histograms (0.00s)
--- PASS: TestProjectFilterReachesEveryLogQuery (0.00s)
=== RUN   TestFilterDataListsProjects
--- PASS: TestFilterDataListsProjects (0.00s)
=== RUN   TestMCPAttributionFilterParsing
--- PASS: TestMCPAttributionFilterParsing (0.00s)
=== RUN   TestLogsStatsWithoutCompareIsUnchanged
--- PASS: TestLogsStatsWithoutCompareIsUnchanged (0.00s)
=== RUN   TestLogsStatsComparePreviousWindow
--- PASS: TestLogsStatsComparePreviousWindow (0.00s)
=== RUN   TestLogsStatsCompareSkippedWhenUnbounded
--- PASS: TestLogsStatsCompareSkippedWhenUnbounded (0.00s)
=== RUN   TestLogsStatsCompareDegradesOnPreviousError
--- PASS: TestLogsStatsCompareDegradesOnPreviousError (0.00s)
=== RUN   TestLogsStatsCompareSkippedForRequestIDLookup
--- PASS: TestLogsStatsCompareSkippedForRequestIDLookup (0.00s)
=== RUN   TestLogsStatsCompareSkippedForEmptyWindow
=== RUN   TestLogsStatsCompareSkippedForEmptyWindow/zero_length
--- PASS: TestLogsStatsCompareSkippedForEmptyWindow/zero_length (0.00s)
=== RUN   TestLogsStatsCompareSkippedForEmptyWindow/reversed
--- PASS: TestLogsStatsCompareSkippedForEmptyWindow/reversed (0.00s)
--- PASS: TestLogsStatsCompareSkippedForEmptyWindow (0.00s)
=== RUN   TestHiddenRequestTypesRoutes
--- PASS: TestHiddenRequestTypesRoutes (0.00s)
=== RUN   TestHiddenRequestTypesFilterCache
--- PASS: TestHiddenRequestTypesFilterCache (0.00s)
=== RUN   TestMCPClientRequestReadsAllowByDefaultUnderBothKeys
=== RUN   TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/current_key_only
--- PASS: TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/current_key_only (0.00s)
=== RUN   TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/earlier_key_only
--- PASS: TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/earlier_key_only (0.00s)
=== RUN   TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/both_sent,_current_wins
--- PASS: TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/both_sent,_current_wins (0.00s)
=== RUN   TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/neither_sent
--- PASS: TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/neither_sent (0.00s)
=== RUN   TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/fields_of_the_outer_request_still_decode
--- PASS: TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/fields_of_the_outer_request_still_decode (0.00s)
--- PASS: TestMCPClientRequestReadsAllowByDefaultUnderBothKeys (0.00s)
=== RUN   TestMCPClientUpdateRequestResolvedAllowByDefault
=== RUN   TestMCPClientUpdateRequestResolvedAllowByDefault/neither_sent_keeps_existing_true
--- PASS: TestMCPClientUpdateRequestResolvedAllowByDefault/neither_sent_keeps_existing_true (0.00s)
=== RUN   TestMCPClientUpdateRequestResolvedAllowByDefault/neither_sent_keeps_existing_false
--- PASS: TestMCPClientUpdateRequestResolvedAllowByDefault/neither_sent_keeps_existing_false (0.00s)
=== RUN   TestMCPClientUpdateRequestResolvedAllowByDefault/current_key_only
--- PASS: TestMCPClientUpdateRequestResolvedAllowByDefault/current_key_only (0.00s)
=== RUN   TestMCPClientUpdateRequestResolvedAllowByDefault/earlier_key_only
--- PASS: TestMCPClientUpdateRequestResolvedAllowByDefault/earlier_key_only (0.00s)
=== RUN   TestMCPClientUpdateRequestResolvedAllowByDefault/both_sent,_current_wins
--- PASS: TestMCPClientUpdateRequestResolvedAllowByDefault/both_sent,_current_wins (0.00s)
--- PASS: TestMCPClientUpdateRequestResolvedAllowByDefault (0.00s)
=== RUN   TestParseOauthScopesJSON
=== RUN   TestParseOauthScopesJSON/empty_string
--- PASS: TestParseOauthScopesJSON/empty_string (0.00s)
=== RUN   TestParseOauthScopesJSON/whitespace_only
--- PASS: TestParseOauthScopesJSON/whitespace_only (0.00s)
=== RUN   TestParseOauthScopesJSON/json_null_(provider_omitted_scope)
--- PASS: TestParseOauthScopesJSON/json_null_(provider_omitted_scope) (0.00s)
=== RUN   TestParseOauthScopesJSON/empty_array
--- PASS: TestParseOauthScopesJSON/empty_array (0.00s)
=== RUN   TestParseOauthScopesJSON/malformed
--- PASS: TestParseOauthScopesJSON/malformed (0.00s)
=== RUN   TestParseOauthScopesJSON/wrong_shape
--- PASS: TestParseOauthScopesJSON/wrong_shape (0.00s)
=== RUN   TestParseOauthScopesJSON/list
--- PASS: TestParseOauthScopesJSON/list (0.00s)
=== RUN   TestParseOauthScopesJSON/drops_blank_entries
--- PASS: TestParseOauthScopesJSON/drops_blank_entries (0.00s)
=== RUN   TestParseOauthScopesJSON/trims_entries
--- PASS: TestParseOauthScopesJSON/trims_entries (0.00s)
--- PASS: TestParseOauthScopesJSON (0.00s)
=== RUN   TestTokenRow_TokenOnlyFields
--- PASS: TestTokenRow_TokenOnlyFields (0.00s)
=== RUN   TestMCPClientCredentialResponse
=== RUN   TestMCPClientCredentialResponse/oauth_reads_the_shared_token,_not_the_admin_one
--- PASS: TestMCPClientCredentialResponse/oauth_reads_the_shared_token,_not_the_admin_one (0.00s)
=== RUN   TestMCPClientCredentialResponse/per_user_oauth_reads_the_admin_token,_not_the_shared_one
--- PASS: TestMCPClientCredentialResponse/per_user_oauth_reads_the_admin_token,_not_the_shared_one (0.00s)
=== RUN   TestMCPClientCredentialResponse/token_exchange_reads_the_admin_token
--- PASS: TestMCPClientCredentialResponse/token_exchange_reads_the_admin_token (0.00s)
=== RUN   TestMCPClientCredentialResponse/per_user_headers_reads_the_admin_header_credential_with_sorted_key_names
--- PASS: TestMCPClientCredentialResponse/per_user_headers_reads_the_admin_header_credential_with_sorted_key_names (0.00s)
=== RUN   TestMCPClientCredentialResponse/none_has_no_self-held_credential
--- PASS: TestMCPClientCredentialResponse/none_has_no_self-held_credential (0.00s)
=== RUN   TestMCPClientCredentialResponse/headers_has_no_self-held_credential
--- PASS: TestMCPClientCredentialResponse/headers_has_no_self-held_credential (0.00s)
=== RUN   TestMCPClientCredentialResponse/missing_rows_yield_no_block
--- PASS: TestMCPClientCredentialResponse/missing_rows_yield_no_block (0.00s)
=== RUN   TestMCPClientCredentialResponse/unreadable_header_values_still_project_status_and_timestamps
--- PASS: TestMCPClientCredentialResponse/unreadable_header_values_still_project_status_and_timestamps (0.00s)
--- PASS: TestMCPClientCredentialResponse (0.00s)
=== RUN   TestUpdateMCPClient_DisabledToEnabled_WithInvalidReplacementHeaders_PreflightRejects
--- PASS: TestUpdateMCPClient_DisabledToEnabled_WithInvalidReplacementHeaders_PreflightRejects (0.00s)
=== RUN   TestEndpointSlugDerivable
--- PASS: TestEndpointSlugDerivable (0.00s)
=== RUN   TestMCPHeadersEqual
=== RUN   TestMCPHeadersEqual/both_empty
--- PASS: TestMCPHeadersEqual/both_empty (0.00s)
=== RUN   TestMCPHeadersEqual/identical_single_header
--- PASS: TestMCPHeadersEqual/identical_single_header (0.00s)
=== RUN   TestMCPHeadersEqual/changed_value
--- PASS: TestMCPHeadersEqual/changed_value (0.00s)
=== RUN   TestMCPHeadersEqual/added_header
--- PASS: TestMCPHeadersEqual/added_header (0.00s)
=== RUN   TestMCPHeadersEqual/removed_header
--- PASS: TestMCPHeadersEqual/removed_header (0.00s)
=== RUN   TestMCPHeadersEqual/renamed_key_with_same_value
--- PASS: TestMCPHeadersEqual/renamed_key_with_same_value (0.00s)
=== RUN   TestMCPHeadersEqual/case-differing_key
--- PASS: TestMCPHeadersEqual/case-differing_key (0.00s)
=== RUN   TestMCPHeadersEqual/empty_vs_populated
--- PASS: TestMCPHeadersEqual/empty_vs_populated (0.00s)
=== RUN   TestMCPHeadersEqual/nil_vs_empty_are_both_no_headers
--- PASS: TestMCPHeadersEqual/nil_vs_empty_are_both_no_headers (0.00s)
--- PASS: TestMCPHeadersEqual (0.00s)
=== RUN   TestIsPrematureOAuthCompletion
=== RUN   TestIsPrematureOAuthCompletion/nil_flow_(already_cleaned_up_by_a_genuine_completion)_is_not_premature
--- PASS: TestIsPrematureOAuthCompletion/nil_flow_(already_cleaned_up_by_a_genuine_completion)_is_not_premature (0.00s)
=== RUN   TestIsPrematureOAuthCompletion/pending,_unexpired_flow_means_the_callback_never_ran:_premature
--- PASS: TestIsPrematureOAuthCompletion/pending,_unexpired_flow_means_the_callback_never_ran:_premature (0.00s)
=== RUN   TestIsPrematureOAuthCompletion/pending_but_expired_flow_(abandoned_attempt)_does_not_block_a_later_fresh_completion
--- PASS: TestIsPrematureOAuthCompletion/pending_but_expired_flow_(abandoned_attempt)_does_not_block_a_later_fresh_completion (0.00s)
=== RUN   TestIsPrematureOAuthCompletion/failed_flow_row_is_not_premature_(already_resolved,_just_unsuccessfully)
--- PASS: TestIsPrematureOAuthCompletion/failed_flow_row_is_not_premature_(already_resolved,_just_unsuccessfully) (0.00s)
--- PASS: TestIsPrematureOAuthCompletion (0.00s)
=== RUN   TestProjectMCPCredentialState
=== RUN   TestProjectMCPCredentialState/per_user_oauth_connected_with_dead_admin_token_projects_needs_reauth
--- PASS: TestProjectMCPCredentialState/per_user_oauth_connected_with_dead_admin_token_projects_needs_reauth (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_oauth_connected_with_active_admin_token_stays_connected
--- PASS: TestProjectMCPCredentialState/per_user_oauth_connected_with_active_admin_token_stays_connected (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_oauth_connected_with_no_admin_token_stays_connected_(pre-retention_client)
--- PASS: TestProjectMCPCredentialState/per_user_oauth_connected_with_no_admin_token_stays_connected_(pre-retention_client) (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_oauth_connected_with_orphaned_admin_token_stays_connected
--- PASS: TestProjectMCPCredentialState/per_user_oauth_connected_with_orphaned_admin_token_stays_connected (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_oauth_ignores_the_header_credential_status
--- PASS: TestProjectMCPCredentialState/per_user_oauth_ignores_the_header_credential_status (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_headers_connected_with_stale_admin_credential_projects_needs_reauth
--- PASS: TestProjectMCPCredentialState/per_user_headers_connected_with_stale_admin_credential_projects_needs_reauth (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_headers_connected_with_active_admin_credential_stays_connected
--- PASS: TestProjectMCPCredentialState/per_user_headers_connected_with_active_admin_credential_stays_connected (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_headers_connected_with_no_admin_credential_stays_connected_(pre-retention_client)
--- PASS: TestProjectMCPCredentialState/per_user_headers_connected_with_no_admin_credential_stays_connected_(pre-retention_client) (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_headers_ignores_the_oauth_token_status
--- PASS: TestProjectMCPCredentialState/per_user_headers_ignores_the_oauth_token_status (0.00s)
=== RUN   TestProjectMCPCredentialState/pending_verification_runtime_state_passes_through_even_with_a_dead_admin_token
--- PASS: TestProjectMCPCredentialState/pending_verification_runtime_state_passes_through_even_with_a_dead_admin_token (0.00s)
=== RUN   TestProjectMCPCredentialState/disabled_runtime_state_passes_through_even_with_a_stale_admin_credential
--- PASS: TestProjectMCPCredentialState/disabled_runtime_state_passes_through_even_with_a_stale_admin_credential (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_oauth_unstable_with_a_dead_admin_token_projects_needs_reauth_(needs_reauth_is_a_bigger_indicator_than_unstable)
--- PASS: TestProjectMCPCredentialState/per_user_oauth_unstable_with_a_dead_admin_token_projects_needs_reauth_(needs_reauth_is_a_bigger_indicator_than_unstable) (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_headers_unstable_with_a_stale_admin_credential_projects_needs_reauth
--- PASS: TestProjectMCPCredentialState/per_user_headers_unstable_with_a_stale_admin_credential_projects_needs_reauth (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_oauth_unstable_with_an_active_admin_token_stays_unstable
--- PASS: TestProjectMCPCredentialState/per_user_oauth_unstable_with_an_active_admin_token_stays_unstable (0.00s)
=== RUN   TestProjectMCPCredentialState/shared_oauth_connected_with_rotated_token_projects_needs_reauth
--- PASS: TestProjectMCPCredentialState/shared_oauth_connected_with_rotated_token_projects_needs_reauth (0.00s)
=== RUN   TestProjectMCPCredentialState/shared_oauth_connected_with_active_token_stays_connected
--- PASS: TestProjectMCPCredentialState/shared_oauth_connected_with_active_token_stays_connected (0.00s)
=== RUN   TestProjectMCPCredentialState/shared_oauth_connected_with_no_token_row_stays_connected
--- PASS: TestProjectMCPCredentialState/shared_oauth_connected_with_no_token_row_stays_connected (0.00s)
=== RUN   TestProjectMCPCredentialState/shared_oauth_unstable_with_rotated_token_projects_needs_reauth
--- PASS: TestProjectMCPCredentialState/shared_oauth_unstable_with_rotated_token_projects_needs_reauth (0.00s)
=== RUN   TestProjectMCPCredentialState/shared_oauth_ignores_the_per-user_admin_credential_statuses
--- PASS: TestProjectMCPCredentialState/shared_oauth_ignores_the_per-user_admin_credential_statuses (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_oauth_ignores_the_shared_token_status
--- PASS: TestProjectMCPCredentialState/per_user_oauth_ignores_the_shared_token_status (0.00s)
=== RUN   TestProjectMCPCredentialState/shared_headers_are_never_projected
--- PASS: TestProjectMCPCredentialState/shared_headers_are_never_projected (0.00s)
=== RUN   TestProjectMCPCredentialState/disabled_shared_oauth_passes_through_even_with_a_rotated_token
--- PASS: TestProjectMCPCredentialState/disabled_shared_oauth_passes_through_even_with_a_rotated_token (0.00s)
=== RUN   TestProjectMCPCredentialState/non-auth_clients_are_never_projected
--- PASS: TestProjectMCPCredentialState/non-auth_clients_are_never_projected (0.00s)
--- PASS: TestProjectMCPCredentialState (0.00s)
=== RUN   TestToolSyncIntervalBoundsRejectNanosecondValues
--- PASS: TestToolSyncIntervalBoundsRejectNanosecondValues (0.00s)
=== RUN   TestToolSyncIntervalBoundsRejectNegatives
--- PASS: TestToolSyncIntervalBoundsRejectNegatives (0.00s)
=== RUN   TestToolSyncIntervalBoundsAcceptRealisticValues
--- PASS: TestToolSyncIntervalBoundsAcceptRealisticValues (0.00s)
=== RUN   TestToolSyncIntervalBoundsSurvivePersistence
--- PASS: TestToolSyncIntervalBoundsSurvivePersistence (0.00s)
=== RUN   TestToolSyncIntervalBoundsPreventOverflow
--- PASS: TestToolSyncIntervalBoundsPreventOverflow (0.00s)
=== RUN   TestUpdateMCPClientCredentialsWithRetry_NotApplicable_ReturnsImmediately
--- PASS: TestUpdateMCPClientCredentialsWithRetry_NotApplicable_ReturnsImmediately (0.00s)
=== RUN   TestUpdateMCPClientCredentialsWithRetry_TransientReconnectError_StillRetries
--- PASS: TestUpdateMCPClientCredentialsWithRetry_TransientReconnectError_StillRetries (1.00s)
=== RUN   TestIsConcurrentVerifyReplay
=== RUN   TestIsConcurrentVerifyReplay/genuine_concurrent_race:_started_pending_verification,_another_request_finished_the_bootstrap_first
--- PASS: TestIsConcurrentVerifyReplay/genuine_concurrent_race:_started_pending_verification,_another_request_finished_the_bootstrap_first (0.00s)
=== RUN   TestIsConcurrentVerifyReplay/repair_on_an_already-verified_client:_not_pending_verification,_tools_already_present
--- PASS: TestIsConcurrentVerifyReplay/repair_on_an_already-verified_client:_not_pending_verification,_tools_already_present (0.00s)
=== RUN   TestIsConcurrentVerifyReplay/first-time_verification_with_no_race:_started_pending_verification,_still_no_tools_after_reload
--- PASS: TestIsConcurrentVerifyReplay/first-time_verification_with_no_race:_started_pending_verification,_still_no_tools_after_reload (0.00s)
--- PASS: TestIsConcurrentVerifyReplay (0.00s)
=== RUN   TestConsentFlowDetail
=== RUN   TestConsentFlowDetail/pending_flow_returns_client_and_modes
--- PASS: TestConsentFlowDetail/pending_flow_returns_client_and_modes (0.00s)
=== RUN   TestConsentFlowDetail/missing_flow_returns_404
--- PASS: TestConsentFlowDetail/missing_flow_returns_404 (0.00s)
=== RUN   TestConsentFlowDetail/empty_flow_id_returns_400
--- PASS: TestConsentFlowDetail/empty_flow_id_returns_400 (0.00s)
=== RUN   TestConsentFlowDetail/expired_flow_returns_410
--- PASS: TestConsentFlowDetail/expired_flow_returns_410 (0.00s)
=== RUN   TestConsentFlowDetail/already-consented_flow_returns_410
--- PASS: TestConsentFlowDetail/already-consented_flow_returns_410 (0.00s)
--- PASS: TestConsentFlowDetail (0.00s)
=== RUN   TestConsentAvailableModes
=== RUN   TestConsentAvailableModes/vk_and_session_when_auth_not_enforced
--- PASS: TestConsentAvailableModes/vk_and_session_when_auth_not_enforced (0.00s)
=== RUN   TestConsentAvailableModes/vk_only_when_auth_enforced
--- PASS: TestConsentAvailableModes/vk_only_when_auth_enforced (0.00s)
=== RUN   TestConsentAvailableModes/adds_user_when_resolver_offers_it
--- PASS: TestConsentAvailableModes/adds_user_when_resolver_offers_it (0.00s)
=== RUN   TestConsentAvailableModes/disable_vk_drops_vk_when_user_mode_available
--- PASS: TestConsentAvailableModes/disable_vk_drops_vk_when_user_mode_available (0.00s)
=== RUN   TestConsentAvailableModes/disable_vk_leaves_user-only_when_auth_enforced
--- PASS: TestConsentAvailableModes/disable_vk_leaves_user-only_when_auth_enforced (0.00s)
=== RUN   TestConsentAvailableModes/disable_vk_ignored_without_user_mode
--- PASS: TestConsentAvailableModes/disable_vk_ignored_without_user_mode (0.00s)
--- PASS: TestConsentAvailableModes (0.00s)
=== RUN   TestConsentFlowSubmit_VK
=== RUN   TestConsentFlowSubmit_VK/active_VK_mints_a_code
--- PASS: TestConsentFlowSubmit_VK/active_VK_mints_a_code (0.00s)
=== RUN   TestConsentFlowSubmit_VK/inactive_VK_is_rejected
--- PASS: TestConsentFlowSubmit_VK/inactive_VK_is_rejected (0.00s)
=== RUN   TestConsentFlowSubmit_VK/unknown_VK_is_rejected
--- PASS: TestConsentFlowSubmit_VK/unknown_VK_is_rejected (0.00s)
=== RUN   TestConsentFlowSubmit_VK/empty_VK_value_is_rejected
--- PASS: TestConsentFlowSubmit_VK/empty_VK_value_is_rejected (0.00s)
=== RUN   TestConsentFlowSubmit_VK/double_submit_returns_410_on_the_second_attempt
--- PASS: TestConsentFlowSubmit_VK/double_submit_returns_410_on_the_second_attempt (0.00s)
--- PASS: TestConsentFlowSubmit_VK (0.00s)
=== RUN   TestConsentFlowSubmit_Session
=== RUN   TestConsentFlowSubmit_Session/session_mode_mints_a_server-side_token_when_auth_not_enforced
--- PASS: TestConsentFlowSubmit_Session/session_mode_mints_a_server-side_token_when_auth_not_enforced (0.00s)
=== RUN   TestConsentFlowSubmit_Session/session_mode_is_unavailable_when_auth_is_enforced
--- PASS: TestConsentFlowSubmit_Session/session_mode_is_unavailable_when_auth_is_enforced (0.00s)
--- PASS: TestConsentFlowSubmit_Session (0.00s)
=== RUN   TestConsentFlowSubmit_User
=== RUN   TestConsentFlowSubmit_User/user_mode_rejected_when_no_resolver_(mode_not_offered)
--- PASS: TestConsentFlowSubmit_User/user_mode_rejected_when_no_resolver_(mode_not_offered) (0.00s)
=== RUN   TestConsentFlowSubmit_User/resolved_session_yields_user_mode
--- PASS: TestConsentFlowSubmit_User/resolved_session_yields_user_mode (0.00s)
=== RUN   TestConsentFlowSubmit_User/user_mode_without_a_session_is_rejected
--- PASS: TestConsentFlowSubmit_User/user_mode_without_a_session_is_rejected (0.00s)
--- PASS: TestConsentFlowSubmit_User (0.00s)
=== RUN   TestConsentFlowSubmit_VKUserBinding
=== RUN   TestConsentFlowSubmit_VKUserBinding/bound_VK_upgrades_to_user_when_logged-in_user_matches
--- PASS: TestConsentFlowSubmit_VKUserBinding/bound_VK_upgrades_to_user_when_logged-in_user_matches (0.00s)
=== RUN   TestConsentFlowSubmit_VKUserBinding/bound_VK_rejected_when_logged-in_user_differs
--- PASS: TestConsentFlowSubmit_VKUserBinding/bound_VK_rejected_when_logged-in_user_differs (0.00s)
=== RUN   TestConsentFlowSubmit_VKUserBinding/bound_VK_rejected_when_not_signed_in
--- PASS: TestConsentFlowSubmit_VKUserBinding/bound_VK_rejected_when_not_signed_in (0.00s)
--- PASS: TestConsentFlowSubmit_VKUserBinding (0.00s)
=== RUN   TestDiscovery_GatedOnAuthMode
=== RUN   TestDiscovery_GatedOnAuthMode/headers_mode_returns_404_on_all_discovery_endpoints
--- PASS: TestDiscovery_GatedOnAuthMode/headers_mode_returns_404_on_all_discovery_endpoints (0.00s)
=== RUN   TestDiscovery_GatedOnAuthMode/oauth_and_both_modes_serve_discovery
--- PASS: TestDiscovery_GatedOnAuthMode/oauth_and_both_modes_serve_discovery (0.00s)
--- PASS: TestDiscovery_GatedOnAuthMode (0.03s)
=== RUN   TestDiscovery_ProtectedResourceMetadata
--- PASS: TestDiscovery_ProtectedResourceMetadata (0.00s)
=== RUN   TestDiscovery_AuthorizationServerMetadata
--- PASS: TestDiscovery_AuthorizationServerMetadata (0.00s)
=== RUN   TestDiscovery_JWKS
--- PASS: TestDiscovery_JWKS (0.02s)
=== RUN   TestHandleRegister_DCR
=== RUN   TestHandleRegister_DCR/valid_registration_returns_201_with_defaults
    mcpoauth2issuance_test.go:123: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:123
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleRegister_DCR/valid_registration_returns_201_with_defaults
--- FAIL: TestHandleRegister_DCR/valid_registration_returns_201_with_defaults (0.00s)
=== RUN   TestHandleRegister_DCR/missing_redirect_uris_is_rejected
    mcpoauth2issuance_test.go:140: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:140
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleRegister_DCR/missing_redirect_uris_is_rejected
--- FAIL: TestHandleRegister_DCR/missing_redirect_uris_is_rejected (0.00s)
=== RUN   TestHandleRegister_DCR/non-public_auth_method_is_rejected
    mcpoauth2issuance_test.go:151: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:151
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleRegister_DCR/non-public_auth_method_is_rejected
--- FAIL: TestHandleRegister_DCR/non-public_auth_method_is_rejected (0.00s)
=== RUN   TestHandleRegister_DCR/malformed_JSON_is_rejected
    mcpoauth2issuance_test.go:162: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:162
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleRegister_DCR/malformed_JSON_is_rejected
--- FAIL: TestHandleRegister_DCR/malformed_JSON_is_rejected (0.00s)
=== RUN   TestHandleRegister_DCR/unsupported_grant_type_is_rejected
    mcpoauth2issuance_test.go:173: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:173
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleRegister_DCR/unsupported_grant_type_is_rejected
--- FAIL: TestHandleRegister_DCR/unsupported_grant_type_is_rejected (0.00s)
=== RUN   TestHandleRegister_DCR/unsupported_response_type_is_rejected
    mcpoauth2issuance_test.go:184: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:184
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleRegister_DCR/unsupported_response_type_is_rejected
--- FAIL: TestHandleRegister_DCR/unsupported_response_type_is_rejected (0.00s)
--- FAIL: TestHandleRegister_DCR (0.01s)
=== RUN   TestHandleAuthorize
=== RUN   TestHandleAuthorize/happy_path_redirects_to_the_consent_page
    mcpoauth2issuance_test.go:209: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:209
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/happy_path_redirects_to_the_consent_page
--- FAIL: TestHandleAuthorize/happy_path_redirects_to_the_consent_page (0.00s)
=== RUN   TestHandleAuthorize/loopback_redirect_matches_on_any_port
    mcpoauth2issuance_test.go:219: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:219
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/loopback_redirect_matches_on_any_port
--- FAIL: TestHandleAuthorize/loopback_redirect_matches_on_any_port (0.00s)
=== RUN   TestHandleAuthorize/unknown_client_id_is_rejected
    mcpoauth2issuance_test.go:229: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:229
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/unknown_client_id_is_rejected
--- FAIL: TestHandleAuthorize/unknown_client_id_is_rejected (0.00s)
=== RUN   TestHandleAuthorize/unregistered_redirect_uri_is_rejected
    mcpoauth2issuance_test.go:238: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:238
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/unregistered_redirect_uri_is_rejected
--- FAIL: TestHandleAuthorize/unregistered_redirect_uri_is_rejected (0.00s)
=== RUN   TestHandleAuthorize/non-code_response_type_redirects_with_error
    mcpoauth2issuance_test.go:260: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:260
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/non-code_response_type_redirects_with_error
--- FAIL: TestHandleAuthorize/non-code_response_type_redirects_with_error (0.00s)
=== RUN   TestHandleAuthorize/non-S256_challenge_method_redirects_with_error
    mcpoauth2issuance_test.go:260: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:260
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/non-S256_challenge_method_redirects_with_error
--- FAIL: TestHandleAuthorize/non-S256_challenge_method_redirects_with_error (0.00s)
=== RUN   TestHandleAuthorize/mismatched_resource_redirects_with_error
    mcpoauth2issuance_test.go:260: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:260
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/mismatched_resource_redirects_with_error
--- FAIL: TestHandleAuthorize/mismatched_resource_redirects_with_error (0.00s)
=== RUN   TestHandleAuthorize/scope_exceeds_registered_redirects_with_error
    mcpoauth2issuance_test.go:260: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:260
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/scope_exceeds_registered_redirects_with_error
--- FAIL: TestHandleAuthorize/scope_exceeds_registered_redirects_with_error (0.00s)
=== RUN   TestHandleAuthorize/omitted_resource_defaults_to_canonical_and_proceeds
    mcpoauth2issuance_test.go:275: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:275
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/omitted_resource_defaults_to_canonical_and_proceeds
--- FAIL: TestHandleAuthorize/omitted_resource_defaults_to_canonical_and_proceeds (0.00s)
--- FAIL: TestHandleAuthorize (0.01s)
=== RUN   TestHandleToken_AuthorizationCode
=== RUN   TestHandleToken_AuthorizationCode/happy_path_issues_a_verifiable_token_pair
    mcpoauth2issuance_test.go:292: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:292
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_AuthorizationCode/happy_path_issues_a_verifiable_token_pair
--- FAIL: TestHandleToken_AuthorizationCode/happy_path_issues_a_verifiable_token_pair (0.00s)
=== RUN   TestHandleToken_AuthorizationCode/PKCE_mismatch_is_rejected
    mcpoauth2issuance_test.go:322: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:322
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_AuthorizationCode/PKCE_mismatch_is_rejected
--- FAIL: TestHandleToken_AuthorizationCode/PKCE_mismatch_is_rejected (0.00s)
=== RUN   TestHandleToken_AuthorizationCode/code_is_single-use
    mcpoauth2issuance_test.go:339: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:339
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_AuthorizationCode/code_is_single-use
--- FAIL: TestHandleToken_AuthorizationCode/code_is_single-use (0.00s)
=== RUN   TestHandleToken_AuthorizationCode/expired_code_is_rejected
    mcpoauth2issuance_test.go:361: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:361
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_AuthorizationCode/expired_code_is_rejected
--- FAIL: TestHandleToken_AuthorizationCode/expired_code_is_rejected (0.00s)
=== RUN   TestHandleToken_AuthorizationCode/client_id_mismatch_is_rejected
    mcpoauth2issuance_test.go:378: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:378
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_AuthorizationCode/client_id_mismatch_is_rejected
--- FAIL: TestHandleToken_AuthorizationCode/client_id_mismatch_is_rejected (0.00s)
=== RUN   TestHandleToken_AuthorizationCode/missing_redirect_uri_is_rejected
    mcpoauth2issuance_test.go:395: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:395
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_AuthorizationCode/missing_redirect_uri_is_rejected
--- FAIL: TestHandleToken_AuthorizationCode/missing_redirect_uri_is_rejected (0.00s)
=== RUN   TestHandleToken_AuthorizationCode/missing_required_fields_are_rejected
    mcpoauth2issuance_test.go:411: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:411
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_AuthorizationCode/missing_required_fields_are_rejected
--- FAIL: TestHandleToken_AuthorizationCode/missing_required_fields_are_rejected (0.00s)
=== RUN   TestHandleToken_AuthorizationCode/unsupported_grant_type_is_rejected
    mcpoauth2issuance_test.go:419: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:419
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_AuthorizationCode/unsupported_grant_type_is_rejected
--- FAIL: TestHandleToken_AuthorizationCode/unsupported_grant_type_is_rejected (0.00s)
--- FAIL: TestHandleToken_AuthorizationCode (0.01s)
=== RUN   TestHandleToken_RefreshRotationAndReplay
=== RUN   TestHandleToken_RefreshRotationAndReplay/rotation_issues_a_new_pair_and_carries_the_family
    mcpoauth2issuance_test.go:441: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:441
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_RefreshRotationAndReplay/rotation_issues_a_new_pair_and_carries_the_family
--- FAIL: TestHandleToken_RefreshRotationAndReplay/rotation_issues_a_new_pair_and_carries_the_family (0.00s)
=== RUN   TestHandleToken_RefreshRotationAndReplay/replaying_a_rotated_token_revokes_the_whole_family
    mcpoauth2issuance_test.go:468: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:468
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_RefreshRotationAndReplay/replaying_a_rotated_token_revokes_the_whole_family
--- FAIL: TestHandleToken_RefreshRotationAndReplay/replaying_a_rotated_token_revokes_the_whole_family (0.00s)
=== RUN   TestHandleToken_RefreshRotationAndReplay/client_id_mismatch_is_rejected
    mcpoauth2issuance_test.go:499: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:499
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_RefreshRotationAndReplay/client_id_mismatch_is_rejected
--- FAIL: TestHandleToken_RefreshRotationAndReplay/client_id_mismatch_is_rejected (0.00s)
=== RUN   TestHandleToken_RefreshRotationAndReplay/missing_fields_are_rejected
    mcpoauth2issuance_test.go:512: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:512
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_RefreshRotationAndReplay/missing_fields_are_rejected
--- FAIL: TestHandleToken_RefreshRotationAndReplay/missing_fields_are_rejected (0.00s)
--- FAIL: TestHandleToken_RefreshRotationAndReplay (0.00s)
=== RUN   TestHandleToken_RefreshVKIdentityDisabled
=== RUN   TestHandleToken_RefreshVKIdentityDisabled/vk_refresh_rejected_when_disabled_and_user_mode_available
    mcpoauth2issuance_test.go:538: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:538
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_RefreshVKIdentityDisabled/vk_refresh_rejected_when_disabled_and_user_mode_available
--- FAIL: TestHandleToken_RefreshVKIdentityDisabled/vk_refresh_rejected_when_disabled_and_user_mode_available (0.00s)
=== RUN   TestHandleToken_RefreshVKIdentityDisabled/vk_refresh_not_blocked_by_flag_when_user_mode_unavailable
    mcpoauth2issuance_test.go:554: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:554
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_RefreshVKIdentityDisabled/vk_refresh_not_blocked_by_flag_when_user_mode_unavailable
--- FAIL: TestHandleToken_RefreshVKIdentityDisabled/vk_refresh_not_blocked_by_flag_when_user_mode_unavailable (0.00s)
--- FAIL: TestHandleToken_RefreshVKIdentityDisabled (0.00s)
=== RUN   TestHandleToken_RefreshUserLiveness
=== RUN   TestHandleToken_RefreshUserLiveness/user_refresh_rejected_when_the_user_is_no_longer_active
    mcpoauth2issuance_test.go:591: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:591
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_RefreshUserLiveness/user_refresh_rejected_when_the_user_is_no_longer_active
--- FAIL: TestHandleToken_RefreshUserLiveness/user_refresh_rejected_when_the_user_is_no_longer_active (0.00s)
=== RUN   TestHandleToken_RefreshUserLiveness/user_refresh_succeeds_when_the_user_is_active
    mcpoauth2issuance_test.go:608: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:608
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_RefreshUserLiveness/user_refresh_succeeds_when_the_user_is_active
--- FAIL: TestHandleToken_RefreshUserLiveness/user_refresh_succeeds_when_the_user_is_active (0.00s)
--- FAIL: TestHandleToken_RefreshUserLiveness (0.00s)
=== RUN   TestUpdateConfig_RejectsAuthCodeTTLAboveMax
=== RUN   TestUpdateConfig_RejectsAuthCodeTTLAboveMax/oauth
    mcpoauth2issuance_test.go:644: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:644
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestUpdateConfig_RejectsAuthCodeTTLAboveMax/oauth
--- FAIL: TestUpdateConfig_RejectsAuthCodeTTLAboveMax/oauth (0.00s)
=== RUN   TestUpdateConfig_RejectsAuthCodeTTLAboveMax/both
    mcpoauth2issuance_test.go:644: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:644
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestUpdateConfig_RejectsAuthCodeTTLAboveMax/both
--- FAIL: TestUpdateConfig_RejectsAuthCodeTTLAboveMax/both (0.00s)
=== RUN   TestUpdateConfig_RejectsAuthCodeTTLAboveMax/headers
    mcpoauth2issuance_test.go:644: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:644
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestUpdateConfig_RejectsAuthCodeTTLAboveMax/headers
--- FAIL: TestUpdateConfig_RejectsAuthCodeTTLAboveMax/headers (0.00s)
--- FAIL: TestUpdateConfig_RejectsAuthCodeTTLAboveMax (0.00s)
=== RUN   TestHandleAuthorize_AuthCodeTTLResolution
=== RUN   TestHandleAuthorize_AuthCodeTTLResolution/above_cap_is_clamped
    mcpoauth2issuance_test.go:676: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:676
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize_AuthCodeTTLResolution/above_cap_is_clamped
--- FAIL: TestHandleAuthorize_AuthCodeTTLResolution/above_cap_is_clamped (0.00s)
=== RUN   TestHandleAuthorize_AuthCodeTTLResolution/zero_falls_back_to_default
    mcpoauth2issuance_test.go:676: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:676
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize_AuthCodeTTLResolution/zero_falls_back_to_default
--- FAIL: TestHandleAuthorize_AuthCodeTTLResolution/zero_falls_back_to_default (0.00s)
=== RUN   TestHandleAuthorize_AuthCodeTTLResolution/in-range_used_verbatim
    mcpoauth2issuance_test.go:676: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:676
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize_AuthCodeTTLResolution/in-range_used_verbatim
--- FAIL: TestHandleAuthorize_AuthCodeTTLResolution/in-range_used_verbatim (0.00s)
--- FAIL: TestHandleAuthorize_AuthCodeTTLResolution (0.00s)
=== RUN   TestExtractBearerJWT
=== RUN   TestExtractBearerJWT/jwt_bearer
--- PASS: TestExtractBearerJWT/jwt_bearer (0.00s)
=== RUN   TestExtractBearerJWT/case-insensitive_scheme
--- PASS: TestExtractBearerJWT/case-insensitive_scheme (0.00s)
=== RUN   TestExtractBearerJWT/virtual_key_is_not_a_jwt
--- PASS: TestExtractBearerJWT/virtual_key_is_not_a_jwt (0.00s)
=== RUN   TestExtractBearerJWT/non-bearer_scheme
--- PASS: TestExtractBearerJWT/non-bearer_scheme (0.00s)
=== RUN   TestExtractBearerJWT/empty_header
--- PASS: TestExtractBearerJWT/empty_header (0.00s)
=== RUN   TestExtractBearerJWT/bearer_without_token
--- PASS: TestExtractBearerJWT/bearer_without_token (0.00s)
--- PASS: TestExtractBearerJWT (0.00s)
=== RUN   TestVerifyMCPJWT_ValidEachMode
=== RUN   TestVerifyMCPJWT_ValidEachMode/user
--- PASS: TestVerifyMCPJWT_ValidEachMode/user (0.00s)
=== RUN   TestVerifyMCPJWT_ValidEachMode/vk
--- PASS: TestVerifyMCPJWT_ValidEachMode/vk (0.00s)
=== RUN   TestVerifyMCPJWT_ValidEachMode/session
--- PASS: TestVerifyMCPJWT_ValidEachMode/session (0.00s)
--- PASS: TestVerifyMCPJWT_ValidEachMode (0.06s)
=== RUN   TestVerifyMCPJWT_Rejections
--- PASS: TestPoolExpiredConnection (1.50s)
=== RUN   TestPoolGetSkipsStaleIdleConnection
--- PASS: TestPoolGetSkipsStaleIdleConnection (0.00s)
=== RUN   TestPoolGetDialsThroughConfiguredHTTPProxy
--- PASS: TestPoolGetDialsThroughConfiguredHTTPProxy (0.00s)
=== RUN   TestPoolGetFailsWithUnreachableProxy
--- PASS: TestPoolGetFailsWithUnreachableProxy (0.00s)
=== RUN   TestSessionManagerCreateAndGet
--- PASS: TestSessionManagerCreateAndGet (0.00s)
=== RUN   TestSessionManagerConnectionLimit
--- PASS: TestSessionManagerConnectionLimit (0.00s)
=== RUN   TestSessionReservationLimitAndCleanup
--- PASS: TestSessionReservationLimitAndCleanup (0.00s)
=== RUN   TestSessionManagerRemove
--- PASS: TestSessionManagerRemove (0.00s)
=== RUN   TestSessionLastResponseID
--- PASS: TestSessionLastResponseID (0.00s)
=== RUN   TestSessionManagerCloseAll
--- PASS: TestSessionManagerCloseAll (0.00s)
=== RUN   TestSessionRealtimeState
--- PASS: TestSessionRealtimeState (0.00s)
=== RUN   TestSessionRecordRealtimeInputUpdatesPendingItemAndIgnoresConsumedLateUpdate
--- PASS: TestSessionRecordRealtimeInputUpdatesPendingItemAndIgnoresConsumedLateUpdate (0.00s)
PASS
ok  	github.com/maximhq/bifrost/transports/bifrost-http/websocket	1.716s
=== RUN   TestVerifyMCPJWT_Rejections/expired_token
--- PASS: TestVerifyMCPJWT_Rejections/expired_token (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/nbf_in_the_future
--- PASS: TestVerifyMCPJWT_Rejections/nbf_in_the_future (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/missing_exp
--- PASS: TestVerifyMCPJWT_Rejections/missing_exp (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/missing_iat
--- PASS: TestVerifyMCPJWT_Rejections/missing_iat (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/issuer_mismatch
--- PASS: TestVerifyMCPJWT_Rejections/issuer_mismatch (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/audience_mismatch
--- PASS: TestVerifyMCPJWT_Rejections/audience_mismatch (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/unknown_kid
--- PASS: TestVerifyMCPJWT_Rejections/unknown_kid (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/wrong_signing_key
--- PASS: TestVerifyMCPJWT_Rejections/wrong_signing_key (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/non-RS256_algorithm_(HS256)
--- PASS: TestVerifyMCPJWT_Rejections/non-RS256_algorithm_(HS256) (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/RS-family_but_not_RS256_(RS384)
--- PASS: TestVerifyMCPJWT_Rejections/RS-family_but_not_RS256_(RS384) (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/alg_none
--- PASS: TestVerifyMCPJWT_Rejections/alg_none (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/malformed_garbage
--- PASS: TestVerifyMCPJWT_Rejections/malformed_garbage (0.00s)
--- PASS: TestVerifyMCPJWT_Rejections (0.21s)
=== RUN   TestVerifyMCPJWT_NilSigningKeyNotLabeledInvalidToken
--- PASS: TestVerifyMCPJWT_NilSigningKeyNotLabeledInvalidToken (0.07s)
=== RUN   TestCachedSigningKey_ConfigFaults
=== RUN   TestCachedSigningKey_ConfigFaults/nil_config_store
--- PASS: TestCachedSigningKey_ConfigFaults/nil_config_store (0.00s)
=== RUN   TestCachedSigningKey_ConfigFaults/signing_key_load_error
--- PASS: TestCachedSigningKey_ConfigFaults/signing_key_load_error (0.00s)
--- PASS: TestCachedSigningKey_ConfigFaults (0.00s)
=== RUN   TestInjectJWTContext
=== RUN   TestInjectJWTContext/user_mode_sets_user_id
--- PASS: TestInjectJWTContext/user_mode_sets_user_id (0.00s)
=== RUN   TestInjectJWTContext/vk_mode_sets_the_raw_vk_value_and_lets_governance_derive_the_id
--- PASS: TestInjectJWTContext/vk_mode_sets_the_raw_vk_value_and_lets_governance_derive_the_id (0.00s)
=== RUN   TestInjectJWTContext/vk_mode_without_vk_errors
--- PASS: TestInjectJWTContext/vk_mode_without_vk_errors (0.00s)
=== RUN   TestInjectJWTContext/session_mode_sets_session_id
--- PASS: TestInjectJWTContext/session_mode_sets_session_id (0.00s)
=== RUN   TestInjectJWTContext/missing_sub_errors
--- PASS: TestInjectJWTContext/missing_sub_errors (0.00s)
=== RUN   TestInjectJWTContext/unknown_bf_mode_errors
--- PASS: TestInjectJWTContext/unknown_bf_mode_errors (0.00s)
--- PASS: TestInjectJWTContext (0.00s)
=== RUN   TestListSessions
=== RUN   TestListSessions/returns_the_grant_rows
--- PASS: TestListSessions/returns_the_grant_rows (0.00s)
=== RUN   TestListSessions/store_error_surfaces_500
--- PASS: TestListSessions/store_error_surfaces_500 (0.00s)
--- PASS: TestListSessions (0.00s)
=== RUN   TestRevokeSession
=== RUN   TestRevokeSession/vk-mode_grant_revokes_without_an_identity_gate
--- PASS: TestRevokeSession/vk-mode_grant_revokes_without_an_identity_gate (0.00s)
=== RUN   TestRevokeSession/user-mode_grant_revokes_when_the_caller_matches_bf_sub
--- PASS: TestRevokeSession/user-mode_grant_revokes_when_the_caller_matches_bf_sub (0.00s)
=== RUN   TestRevokeSession/revokes_a_visible_grant_regardless_of_caller_identity
--- PASS: TestRevokeSession/revokes_a_visible_grant_regardless_of_caller_identity (0.00s)
=== RUN   TestRevokeSession/not-visible_grant_returns_404_without_attempting_revoke
--- PASS: TestRevokeSession/not-visible_grant_returns_404_without_attempting_revoke (0.00s)
=== RUN   TestRevokeSession/empty_id_returns_400
--- PASS: TestRevokeSession/empty_id_returns_400 (0.00s)
--- PASS: TestRevokeSession (0.00s)
=== RUN   TestMatchRedirectURI
=== RUN   TestMatchRedirectURI/exact_non-loopback_match
--- PASS: TestMatchRedirectURI/exact_non-loopback_match (0.00s)
=== RUN   TestMatchRedirectURI/non-loopback_mismatch
--- PASS: TestMatchRedirectURI/non-loopback_mismatch (0.00s)
=== RUN   TestMatchRedirectURI/loopback_any_port_(127.0.0.1)
--- PASS: TestMatchRedirectURI/loopback_any_port_(127.0.0.1) (0.00s)
=== RUN   TestMatchRedirectURI/loopback_any_port_(localhost)
--- PASS: TestMatchRedirectURI/loopback_any_port_(localhost) (0.00s)
=== RUN   TestMatchRedirectURI/loopback_path_must_still_match
--- PASS: TestMatchRedirectURI/loopback_path_must_still_match (0.00s)
=== RUN   TestMatchRedirectURI/loopback_scheme_must_still_match
--- PASS: TestMatchRedirectURI/loopback_scheme_must_still_match (0.00s)
=== RUN   TestMatchRedirectURI/malformed_candidate
--- PASS: TestMatchRedirectURI/malformed_candidate (0.00s)
=== RUN   TestMatchRedirectURI/no_registered_uris
--- PASS: TestMatchRedirectURI/no_registered_uris (0.00s)
--- PASS: TestMatchRedirectURI (0.00s)
=== RUN   TestOAuth2IssuerURL
=== RUN   TestOAuth2IssuerURL/uses_the_configured_issuer_when_set
--- PASS: TestOAuth2IssuerURL/uses_the_configured_issuer_when_set (0.00s)
=== RUN   TestOAuth2IssuerURL/falls_back_to_the_request_host_when_issuer_is_unset
--- PASS: TestOAuth2IssuerURL/falls_back_to_the_request_host_when_issuer_is_unset (0.00s)
--- PASS: TestOAuth2IssuerURL (0.00s)
=== RUN   TestOAuth2ServerCfg_DefaultsWhenUnset
--- PASS: TestOAuth2ServerCfg_DefaultsWhenUnset (0.00s)
=== RUN   TestRejectStdioMCPClientIfAuthBypassed_UnauthenticatedRejected
--- PASS: TestRejectStdioMCPClientIfAuthBypassed_UnauthenticatedRejected (0.00s)
=== RUN   TestRejectStdioMCPClientIfAuthBypassed_AuthenticatedAllowed
--- PASS: TestRejectStdioMCPClientIfAuthBypassed_AuthenticatedAllowed (0.00s)
=== RUN   TestRejectStdioMCPClientIfAuthBypassed_HTTPUnaffected
--- PASS: TestRejectStdioMCPClientIfAuthBypassed_HTTPUnaffected (0.00s)
=== RUN   TestRejectPrivateMCPTargetIfAuthBypassed_UnauthenticatedLoopbackRejected
--- PASS: TestRejectPrivateMCPTargetIfAuthBypassed_UnauthenticatedLoopbackRejected (0.00s)
=== RUN   TestRejectPrivateMCPTargetIfAuthBypassed_UnauthenticatedLinkLocalRejected
--- PASS: TestRejectPrivateMCPTargetIfAuthBypassed_UnauthenticatedLinkLocalRejected (0.00s)
=== RUN   TestRejectPrivateMCPTargetIfAuthBypassed_AuthenticatedLoopbackAllowed
--- PASS: TestRejectPrivateMCPTargetIfAuthBypassed_AuthenticatedLoopbackAllowed (0.00s)
=== RUN   TestRejectPrivateMCPTargetIfAuthBypassed_UnauthenticatedPublicTargetAllowed
--- PASS: TestRejectPrivateMCPTargetIfAuthBypassed_UnauthenticatedPublicTargetAllowed (0.00s)
=== RUN   TestRejectPrivateMCPTargetIfAuthBypassed_STDIOUnaffected
--- PASS: TestRejectPrivateMCPTargetIfAuthBypassed_STDIOUnaffected (0.00s)
=== RUN   TestGetVKFromRequest
=== RUN   TestGetVKFromRequest/x-bf-vk_header
--- PASS: TestGetVKFromRequest/x-bf-vk_header (0.00s)
=== RUN   TestGetVKFromRequest/Authorization_Bearer_header
--- PASS: TestGetVKFromRequest/Authorization_Bearer_header (0.00s)
=== RUN   TestGetVKFromRequest/x-api-key_header
--- PASS: TestGetVKFromRequest/x-api-key_header (0.00s)
=== RUN   TestGetVKFromRequest/x-goog-api-key_header
--- PASS: TestGetVKFromRequest/x-goog-api-key_header (0.00s)
=== RUN   TestGetVKFromRequest/no_header_returns_empty_string
--- PASS: TestGetVKFromRequest/no_header_returns_empty_string (0.00s)
=== RUN   TestGetVKFromRequest/non-VK_Bearer_token_returns_empty_string
--- PASS: TestGetVKFromRequest/non-VK_Bearer_token_returns_empty_string (0.00s)
=== RUN   TestGetVKFromRequest/non-VK_x-goog-api-key_returns_empty_string
--- PASS: TestGetVKFromRequest/non-VK_x-goog-api-key_returns_empty_string (0.00s)
=== RUN   TestGetVKFromRequest/x-bf-vk_takes_priority_over_x-goog-api-key
--- PASS: TestGetVKFromRequest/x-bf-vk_takes_priority_over_x-goog-api-key (0.00s)
--- PASS: TestGetVKFromRequest (0.00s)
=== RUN   TestAuthenticate_JWTPath
=== RUN   TestAuthenticate_JWTPath/oauth_mode:_vk_JWT_stamps_the_key's_value_for_governance_to_resolve
--- PASS: TestAuthenticate_JWTPath/oauth_mode:_vk_JWT_stamps_the_key's_value_for_governance_to_resolve (0.00s)
=== RUN   TestAuthenticate_JWTPath/vk_JWT_with_an_inactive_key_is_stamped_rather_than_refused_here
--- PASS: TestAuthenticate_JWTPath/vk_JWT_with_an_inactive_key_is_stamped_rather_than_refused_here (0.00s)
=== RUN   TestAuthenticate_JWTPath/vk_JWT_for_an_unknown_key_is_rejected
--- PASS: TestAuthenticate_JWTPath/vk_JWT_for_an_unknown_key_is_rejected (0.00s)
=== RUN   TestAuthenticate_JWTPath/vk_JWT_is_rejected_once_virtual-key_identity_is_disabled_and_user_mode_is_offered
--- PASS: TestAuthenticate_JWTPath/vk_JWT_is_rejected_once_virtual-key_identity_is_disabled_and_user_mode_is_offered (0.00s)
=== RUN   TestAuthenticate_JWTPath/user_JWT_without_a_session_stamps_the_user
--- PASS: TestAuthenticate_JWTPath/user_JWT_without_a_session_stamps_the_user (0.00s)
=== RUN   TestAuthenticate_JWTPath/user_JWT_is_rejected_when_the_user_is_no_longer_active
--- PASS: TestAuthenticate_JWTPath/user_JWT_is_rejected_when_the_user_is_no_longer_active (0.00s)
=== RUN   TestAuthenticate_JWTPath/user_JWT_with_a_matching_session_is_accepted
--- PASS: TestAuthenticate_JWTPath/user_JWT_with_a_matching_session_is_accepted (0.00s)
=== RUN   TestAuthenticate_JWTPath/user_JWT_with_a_mismatched_session_is_rejected
--- PASS: TestAuthenticate_JWTPath/user_JWT_with_a_mismatched_session_is_rejected (0.00s)
=== RUN   TestAuthenticate_JWTPath/both_mode:_a_bearer_an_upstream_auth_layer_authenticated_is_accepted_as_the_stamped_user
--- PASS: TestAuthenticate_JWTPath/both_mode:_a_bearer_an_upstream_auth_layer_authenticated_is_accepted_as_the_stamped_user (0.00s)
=== RUN   TestAuthenticate_JWTPath/both_mode:_a_bearer_no_upstream_auth_layer_authenticated_is_refused_on_its_key_id
--- PASS: TestAuthenticate_JWTPath/both_mode:_a_bearer_no_upstream_auth_layer_authenticated_is_refused_on_its_key_id (0.00s)
=== RUN   TestAuthenticate_JWTPath/oauth_strict_mode_verifies_the_bearer_even_when_an_identity_is_stamped_upstream
--- PASS: TestAuthenticate_JWTPath/oauth_strict_mode_verifies_the_bearer_even_when_an_identity_is_stamped_upstream (0.00s)
=== RUN   TestAuthenticate_JWTPath/session_JWT_is_rejected_when_auth_is_enforced
--- PASS: TestAuthenticate_JWTPath/session_JWT_is_rejected_when_auth_is_enforced (0.00s)
=== RUN   TestAuthenticate_JWTPath/session_JWT_stamps_the_session_when_auth_is_not_enforced
--- PASS: TestAuthenticate_JWTPath/session_JWT_stamps_the_session_when_auth_is_not_enforced (0.00s)
=== RUN   TestAuthenticate_JWTPath/both_mode:_session_token_with_a_header_VK_is_rejected
--- PASS: TestAuthenticate_JWTPath/both_mode:_session_token_with_a_header_VK_is_rejected (0.00s)
=== RUN   TestAuthenticate_JWTPath/both_mode:_vk_token_with_a_header_VK_is_rejected
--- PASS: TestAuthenticate_JWTPath/both_mode:_vk_token_with_a_header_VK_is_rejected (0.00s)
--- PASS: TestAuthenticate_JWTPath (0.05s)
=== RUN   TestAuthenticate_HeaderAndAnonPath
=== RUN   TestAuthenticate_HeaderAndAnonPath/headers_mode:_a_header_VK_is_accepted_without_being_looked_up
--- PASS: TestAuthenticate_HeaderAndAnonPath/headers_mode:_a_header_VK_is_accepted_without_being_looked_up (0.00s)
=== RUN   TestAuthenticate_HeaderAndAnonPath/anonymous_access_is_accepted_when_auth_is_not_enforced
--- PASS: TestAuthenticate_HeaderAndAnonPath/anonymous_access_is_accepted_when_auth_is_not_enforced (0.00s)
=== RUN   TestAuthenticate_HeaderAndAnonPath/no_credentials_rejected_when_auth_is_enforced
--- PASS: TestAuthenticate_HeaderAndAnonPath/no_credentials_rejected_when_auth_is_enforced (0.00s)
=== RUN   TestAuthenticate_HeaderAndAnonPath/headers_mode:_an_identity_stamped_upstream_satisfies_enforced_auth
--- PASS: TestAuthenticate_HeaderAndAnonPath/headers_mode:_an_identity_stamped_upstream_satisfies_enforced_auth (0.00s)
=== RUN   TestAuthenticate_HeaderAndAnonPath/both_mode:_an_identity_stamped_upstream_satisfies_enforced_auth
--- PASS: TestAuthenticate_HeaderAndAnonPath/both_mode:_an_identity_stamped_upstream_satisfies_enforced_auth (0.00s)
=== RUN   TestAuthenticate_HeaderAndAnonPath/oauth_strict_mode_ignores_an_identity_stamped_upstream
--- PASS: TestAuthenticate_HeaderAndAnonPath/oauth_strict_mode_ignores_an_identity_stamped_upstream (0.00s)
=== RUN   TestAuthenticate_HeaderAndAnonPath/oauth_strict_mode_rejects_a_header_VK_with_WWW-Authenticate
--- PASS: TestAuthenticate_HeaderAndAnonPath/oauth_strict_mode_rejects_a_header_VK_with_WWW-Authenticate (0.00s)
=== RUN   TestAuthenticate_HeaderAndAnonPath/oauth_strict_mode_with_no_credentials_sets_WWW-Authenticate
--- PASS: TestAuthenticate_HeaderAndAnonPath/oauth_strict_mode_with_no_credentials_sets_WWW-Authenticate (0.00s)
=== RUN   TestAuthenticate_HeaderAndAnonPath/headers_mode:_a_JWT_bearer_is_not_treated_as_a_credential
--- PASS: TestAuthenticate_HeaderAndAnonPath/headers_mode:_a_JWT_bearer_is_not_treated_as_a_credential (0.00s)
--- PASS: TestAuthenticate_HeaderAndAnonPath (0.03s)
=== RUN   TestAdmit
=== RUN   TestAdmit/a_refusal_is_answered_with_the_status_it_carries
--- PASS: TestAdmit/a_refusal_is_answered_with_the_status_it_carries (0.00s)
=== RUN   TestAdmit/a_401_refusal_points_at_the_authorization_server_when_discovery_is_enabled
--- PASS: TestAdmit/a_401_refusal_points_at_the_authorization_server_when_discovery_is_enabled (0.00s)
=== RUN   TestAdmit/a_refusal_without_a_status_is_forbidden
--- PASS: TestAdmit/a_refusal_without_a_status_is_forbidden (0.00s)
=== RUN   TestAdmit/restricted_access_stamps_the_tools_it_grants
--- PASS: TestAdmit/restricted_access_stamps_the_tools_it_grants (0.00s)
=== RUN   TestAdmit/access_granting_no_tool_stamps_an_empty_list,_which_permits_none
--- PASS: TestAdmit/access_granting_no_tool_stamps_an_empty_list,_which_permits_none (0.00s)
=== RUN   TestAdmit/a_caller's_list_narrows_within_the_access_and_cannot_widen_it
--- PASS: TestAdmit/a_caller's_list_narrows_within_the_access_and_cannot_widen_it (0.00s)
=== RUN   TestAdmit/no_access_leaves_the_caller's_list_alone_and_stamps_nothing_of_its_own
--- PASS: TestAdmit/no_access_leaves_the_caller's_list_alone_and_stamps_nothing_of_its_own (0.00s)
=== RUN   TestAdmit/nothing_resolved_stamps_nothing
--- PASS: TestAdmit/nothing_resolved_stamps_nothing (0.00s)
=== RUN   TestAdmit/governance_is_asked_about_the_request's_own_context,_after_authentication
--- PASS: TestAdmit/governance_is_asked_about_the_request's_own_context,_after_authentication (0.00s)
=== RUN   TestAdmit/a_request_that_fails_authentication_is_refused_before_governance_is_asked
--- PASS: TestAdmit/a_request_that_fails_authentication_is_refused_before_governance_is_asked (0.00s)
--- PASS: TestAdmit (0.05s)
=== RUN   TestConvertToolFunctionParametersToMCPInputSchemaPreservesDefs
--- PASS: TestConvertToolFunctionParametersToMCPInputSchemaPreservesDefs (0.00s)
=== RUN   TestConvertToolFunctionParametersToMCPInputSchemaPreservesLegacyDefinitionsAsDefs
--- PASS: TestConvertToolFunctionParametersToMCPInputSchemaPreservesLegacyDefinitionsAsDefs (0.00s)
=== RUN   TestParseMCPSessionsListQuery_Defaults
--- PASS: TestParseMCPSessionsListQuery_Defaults (0.00s)
=== RUN   TestParseMCPSessionsListQuery_AllParams
--- PASS: TestParseMCPSessionsListQuery_AllParams (0.00s)
=== RUN   TestParseMCPSessionsListQuery_LimitCappedAtMax
--- PASS: TestParseMCPSessionsListQuery_LimitCappedAtMax (0.00s)
=== RUN   TestParseMCPSessionsListQuery_InvalidLimit
--- PASS: TestParseMCPSessionsListQuery_InvalidLimit (0.00s)
=== RUN   TestParseMCPSessionsListQuery_InvalidOffset
--- PASS: TestParseMCPSessionsListQuery_InvalidOffset (0.00s)
=== RUN   TestMCPSessionsListQuery_KindAllowed
--- PASS: TestMCPSessionsListQuery_KindAllowed (0.00s)
=== RUN   TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders
=== RUN   TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//anthropic_passthrough/v1/messages/false
--- PASS: TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//anthropic_passthrough/v1/messages/false (0.00s)
=== RUN   TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//anthropic_passthrough/v1/messages/true
--- PASS: TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//anthropic_passthrough/v1/messages/true (0.00s)
=== RUN   TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//chatgpt_passthrough/backend-api/codex/responses/false
--- PASS: TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//chatgpt_passthrough/backend-api/codex/responses/false (0.00s)
=== RUN   TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//chatgpt_passthrough/backend-api/codex/responses/true
--- PASS: TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//chatgpt_passthrough/backend-api/codex/responses/true (0.00s)
=== RUN   TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//openai_passthrough/v1/responses/false
--- PASS: TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//openai_passthrough/v1/responses/false (0.00s)
=== RUN   TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//openai_passthrough/v1/responses/true
--- PASS: TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//openai_passthrough/v1/responses/true (0.00s)
--- PASS: TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders (0.00s)
=== RUN   TestCorsMiddleware_LocalhostOrigins
=== RUN   TestCorsMiddleware_LocalhostOrigins/http://localhost:3000
--- PASS: TestCorsMiddleware_LocalhostOrigins/http://localhost:3000 (0.00s)
=== RUN   TestCorsMiddleware_LocalhostOrigins/https://localhost:3000
--- PASS: TestCorsMiddleware_LocalhostOrigins/https://localhost:3000 (0.00s)
=== RUN   TestCorsMiddleware_LocalhostOrigins/http://127.0.0.1:8080
--- PASS: TestCorsMiddleware_LocalhostOrigins/http://127.0.0.1:8080 (0.00s)
=== RUN   TestCorsMiddleware_LocalhostOrigins/http://0.0.0.0:5000
--- PASS: TestCorsMiddleware_LocalhostOrigins/http://0.0.0.0:5000 (0.00s)
=== RUN   TestCorsMiddleware_LocalhostOrigins/https://127.0.0.1:3000
--- PASS: TestCorsMiddleware_LocalhostOrigins/https://127.0.0.1:3000 (0.00s)
--- PASS: TestCorsMiddleware_LocalhostOrigins (0.00s)
=== RUN   TestCorsMiddleware_ConfiguredOrigins
--- PASS: TestCorsMiddleware_ConfiguredOrigins (0.00s)
=== RUN   TestCorsMiddleware_NonAllowedOrigins
--- PASS: TestCorsMiddleware_NonAllowedOrigins (0.00s)
=== RUN   TestCorsMiddleware_PreflightAllowedOrigin
--- PASS: TestCorsMiddleware_PreflightAllowedOrigin (0.00s)
=== RUN   TestCorsMiddleware_PreflightNonAllowedOrigin
--- PASS: TestCorsMiddleware_PreflightNonAllowedOrigin (0.00s)
=== RUN   TestCorsMiddleware_PreflightLocalhost
--- PASS: TestCorsMiddleware_PreflightLocalhost (0.00s)
=== RUN   TestCorsMiddleware_NoOriginHeader
--- PASS: TestCorsMiddleware_NoOriginHeader (0.00s)
=== RUN   TestChainMiddlewares_NoMiddlewares
--- PASS: TestChainMiddlewares_NoMiddlewares (0.00s)
=== RUN   TestChainMiddlewares_SingleMiddleware
--- PASS: TestChainMiddlewares_SingleMiddleware (0.00s)
=== RUN   TestChainMiddlewares_MultipleMiddlewares
--- PASS: TestChainMiddlewares_MultipleMiddlewares (0.00s)
=== RUN   TestChainMiddlewares_MiddlewareCanModifyContext
--- PASS: TestChainMiddlewares_MiddlewareCanModifyContext (0.00s)
=== RUN   TestIsInferenceWSEndpoint
--- PASS: TestIsInferenceWSEndpoint (0.00s)
=== RUN   TestIsRealtimeTransportEndpoint
--- PASS: TestIsRealtimeTransportEndpoint (0.00s)
=== RUN   TestChainMiddlewares_ShortCircuit
--- PASS: TestChainMiddlewares_ShortCircuit (0.00s)
=== RUN   TestChainMiddlewares_ShortCircuitMiddlePosition
--- PASS: TestChainMiddlewares_ShortCircuitMiddlePosition (0.00s)
=== RUN   TestAuthMiddleware_NilAuthConfig
--- PASS: TestAuthMiddleware_NilAuthConfig (0.00s)
=== RUN   TestAuthMiddleware_DisabledAuthConfig
--- PASS: TestAuthMiddleware_DisabledAuthConfig (0.00s)
=== RUN   TestAuthMiddleware_EnabledAuthConfig_NoAuth
--- PASS: TestAuthMiddleware_EnabledAuthConfig_NoAuth (0.00s)
=== RUN   TestAuthMiddleware_SkillsPublicServeManagementSplit
=== RUN   TestAuthMiddleware_SkillsPublicServeManagementSplit/serve_routes_bypass_auth
--- PASS: TestAuthMiddleware_SkillsPublicServeManagementSplit/serve_routes_bypass_auth (0.00s)
=== RUN   TestAuthMiddleware_SkillsPublicServeManagementSplit/management_routes_require_auth
--- PASS: TestAuthMiddleware_SkillsPublicServeManagementSplit/management_routes_require_auth (0.00s)
--- PASS: TestAuthMiddleware_SkillsPublicServeManagementSplit (0.00s)
=== RUN   TestAuthMiddleware_WhitelistedRoutes
=== RUN   TestAuthMiddleware_WhitelistedRoutes//api/session/is-auth-enabled
--- PASS: TestAuthMiddleware_WhitelistedRoutes//api/session/is-auth-enabled (0.00s)
=== RUN   TestAuthMiddleware_WhitelistedRoutes//api/session/login
--- PASS: TestAuthMiddleware_WhitelistedRoutes//api/session/login (0.00s)
=== RUN   TestAuthMiddleware_WhitelistedRoutes//api/session/logout
--- PASS: TestAuthMiddleware_WhitelistedRoutes//api/session/logout (0.00s)
=== RUN   TestAuthMiddleware_WhitelistedRoutes//api/oauth/callback
--- PASS: TestAuthMiddleware_WhitelistedRoutes//api/oauth/callback (0.00s)
=== RUN   TestAuthMiddleware_WhitelistedRoutes//health
--- PASS: TestAuthMiddleware_WhitelistedRoutes//health (0.00s)
--- PASS: TestAuthMiddleware_WhitelistedRoutes (0.00s)
=== RUN   TestAuthMiddleware_APIMiddleware_DevPrefixDoesNotMatchDevices
=== RUN   TestAuthMiddleware_APIMiddleware_DevPrefixDoesNotMatchDevices/dev_pprof_is_whitelisted
--- PASS: TestAuthMiddleware_APIMiddleware_DevPrefixDoesNotMatchDevices/dev_pprof_is_whitelisted (0.00s)
=== RUN   TestAuthMiddleware_APIMiddleware_DevPrefixDoesNotMatchDevices/dev_pprof_subpath_is_whitelisted
--- PASS: TestAuthMiddleware_APIMiddleware_DevPrefixDoesNotMatchDevices/dev_pprof_subpath_is_whitelisted (0.00s)
=== RUN   TestAuthMiddleware_APIMiddleware_DevPrefixDoesNotMatchDevices/devices_is_NOT_whitelisted
--- PASS: TestAuthMiddleware_APIMiddleware_DevPrefixDoesNotMatchDevices/devices_is_NOT_whitelisted (0.00s)
--- PASS: TestAuthMiddleware_APIMiddleware_DevPrefixDoesNotMatchDevices (0.00s)
=== RUN   TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth
=== RUN   TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth//v1/realtime
--- PASS: TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth//v1/realtime (0.00s)
=== RUN   TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth//openai/v1/realtime
--- PASS: TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth//openai/v1/realtime (0.00s)
=== RUN   TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth//v1/realtime/calls?model=gpt-realtime
--- PASS: TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth//v1/realtime/calls?model=gpt-realtime (0.00s)
=== RUN   TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth//openai/v1/realtime/calls?model=gpt-realtime
--- PASS: TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth//openai/v1/realtime/calls?model=gpt-realtime (0.00s)
--- PASS: TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth (0.00s)
=== RUN   TestAuthMiddleware_InferenceMiddleware_DelegatesAuthToGovernance
=== RUN   TestAuthMiddleware_InferenceMiddleware_DelegatesAuthToGovernance/chat_completion_with_virtual_key
--- PASS: TestAuthMiddleware_InferenceMiddleware_DelegatesAuthToGovernance/chat_completion_with_virtual_key (0.00s)
=== RUN   TestAuthMiddleware_InferenceMiddleware_DelegatesAuthToGovernance/chat_completion_without_credentials
--- PASS: TestAuthMiddleware_InferenceMiddleware_DelegatesAuthToGovernance/chat_completion_without_credentials (0.00s)
=== RUN   TestAuthMiddleware_InferenceMiddleware_DelegatesAuthToGovernance/realtime_minting_(client_secrets)
--- PASS: TestAuthMiddleware_InferenceMiddleware_DelegatesAuthToGovernance/realtime_minting_(client_secrets) (0.00s)
--- PASS: TestAuthMiddleware_InferenceMiddleware_DelegatesAuthToGovernance (0.00s)
=== RUN   TestAuthMiddleware_APIMiddleware_VirtualKeyDoesNotBypass
--- PASS: TestAuthMiddleware_APIMiddleware_VirtualKeyDoesNotBypass (0.00s)
=== RUN   TestAuthMiddleware_UpdateAuthConfig_NilToEnabled
--- PASS: TestAuthMiddleware_UpdateAuthConfig_NilToEnabled (0.00s)
=== RUN   TestAuthMiddleware_BootstrapToken_NoAdminNoToken
--- PASS: TestAuthMiddleware_BootstrapToken_NoAdminNoToken (0.00s)
=== RUN   TestAuthMiddleware_BootstrapToken_AdminExists
--- PASS: TestAuthMiddleware_BootstrapToken_AdminExists (0.00s)
=== RUN   TestAuthMiddleware_BootstrapToken_ValidatesAndClears
--- PASS: TestAuthMiddleware_BootstrapToken_ValidatesAndClears (0.00s)
=== RUN   TestAuthMiddleware_UpdateAuthConfig_EnabledToDisabled
--- PASS: TestAuthMiddleware_UpdateAuthConfig_EnabledToDisabled (0.00s)
=== RUN   TestFasthttpToHTTPRequest
--- PASS: TestFasthttpToHTTPRequest (0.00s)
=== RUN   TestCorsMiddleware_DefaultHeaders
--- PASS: TestCorsMiddleware_DefaultHeaders (0.00s)
=== RUN   TestCorsMiddleware_WildcardHeaders_NonCredentialed
--- PASS: TestCorsMiddleware_WildcardHeaders_NonCredentialed (0.00s)
=== RUN   TestCorsMiddleware_WildcardHeaders_CredentialedPreflight
--- PASS: TestCorsMiddleware_WildcardHeaders_CredentialedPreflight (0.00s)
=== RUN   TestCorsMiddleware_WildcardHeaders_CredentialedNonPreflight
--- PASS: TestCorsMiddleware_WildcardHeaders_CredentialedNonPreflight (0.00s)
=== RUN   TestCorsMiddleware_CustomHeaders
--- PASS: TestCorsMiddleware_CustomHeaders (0.00s)
=== RUN   TestCorsMiddleware_DuplicateHeaders
--- PASS: TestCorsMiddleware_DuplicateHeaders (0.00s)
=== RUN   TestCorsMiddleware_CustomHeadersWithLocalhost
--- PASS: TestCorsMiddleware_CustomHeadersWithLocalhost (0.00s)
=== RUN   TestCorsMiddleware_CustomHeadersNotSetForNonAllowedOrigin
--- PASS: TestCorsMiddleware_CustomHeadersNotSetForNonAllowedOrigin (0.00s)
=== RUN   TestFasthttpToHTTPRequest_PathParams
--- PASS: TestFasthttpToHTTPRequest_PathParams (0.00s)
=== RUN   TestRequestDecompressionMiddleware_SupportedEncodings
=== RUN   TestRequestDecompressionMiddleware_SupportedEncodings/gzip
--- PASS: TestRequestDecompressionMiddleware_SupportedEncodings/gzip (0.00s)
=== RUN   TestRequestDecompressionMiddleware_SupportedEncodings/deflate
--- PASS: TestRequestDecompressionMiddleware_SupportedEncodings/deflate (0.00s)
=== RUN   TestRequestDecompressionMiddleware_SupportedEncodings/brotli
--- PASS: TestRequestDecompressionMiddleware_SupportedEncodings/brotli (0.00s)
=== RUN   TestRequestDecompressionMiddleware_SupportedEncodings/zstd
--- PASS: TestRequestDecompressionMiddleware_SupportedEncodings/zstd (0.00s)
--- PASS: TestRequestDecompressionMiddleware_SupportedEncodings (0.00s)
=== RUN   TestRequestDecompressionMiddleware_InvalidCompressedBody
--- PASS: TestRequestDecompressionMiddleware_InvalidCompressedBody (0.00s)
=== RUN   TestRequestDecompressionMiddleware_UnsupportedEncoding
--- PASS: TestRequestDecompressionMiddleware_UnsupportedEncoding (0.00s)
=== RUN   TestRequestDecompressionMiddleware_DecompressedSizeLimit
--- PASS: TestRequestDecompressionMiddleware_DecompressedSizeLimit (0.00s)
=== RUN   TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding
=== RUN   TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding/gzip
--- PASS: TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding/gzip (0.00s)
=== RUN   TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding/deflate
--- PASS: TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding/deflate (0.00s)
=== RUN   TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding/br
--- PASS: TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding/br (0.00s)
=== RUN   TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding/zstd
--- PASS: TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding/zstd (0.00s)
--- PASS: TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding (0.00s)
=== RUN   TestRequestDecompressionMiddleware_NoContentEncoding
--- PASS: TestRequestDecompressionMiddleware_NoContentEncoding (0.00s)
=== RUN   TestRequestDecompressionMiddleware_ExactSizeLimit
--- PASS: TestRequestDecompressionMiddleware_ExactSizeLimit (0.00s)
=== RUN   TestShouldStreamDecompress
=== RUN   TestShouldStreamDecompress/chunked_(CL=-1)
--- PASS: TestShouldStreamDecompress/chunked_(CL=-1) (0.00s)
=== RUN   TestShouldStreamDecompress/empty_body_(CL=0)
--- PASS: TestShouldStreamDecompress/empty_body_(CL=0) (0.00s)
=== RUN   TestShouldStreamDecompress/small_body
--- PASS: TestShouldStreamDecompress/small_body (0.00s)
=== RUN   TestShouldStreamDecompress/at_default_threshold
--- PASS: TestShouldStreamDecompress/at_default_threshold (0.00s)
=== RUN   TestShouldStreamDecompress/above_default_threshold
--- PASS: TestShouldStreamDecompress/above_default_threshold (0.00s)
=== RUN   TestShouldStreamDecompress/above_custom_threshold
--- PASS: TestShouldStreamDecompress/above_custom_threshold (0.00s)
=== RUN   TestShouldStreamDecompress/below_custom_threshold
--- PASS: TestShouldStreamDecompress/below_custom_threshold (0.00s)
=== RUN   TestShouldStreamDecompress/chunked_with_custom_threshold
--- PASS: TestShouldStreamDecompress/chunked_with_custom_threshold (0.00s)
--- PASS: TestShouldStreamDecompress (0.00s)
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_ChunkedGzip
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_ChunkedGzip (0.00s)
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_AllEncodings
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_AllEncodings/gzip
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_AllEncodings/gzip (0.00s)
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_AllEncodings/deflate
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_AllEncodings/deflate (0.00s)
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_AllEncodings/brotli
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_AllEncodings/brotli (0.00s)
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_AllEncodings/zstd
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_AllEncodings/zstd (0.00s)
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_AllEncodings (0.00s)
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_InvalidBody
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_InvalidBody (0.00s)
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_UnsupportedEncoding
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_UnsupportedEncoding (0.00s)
=== RUN   TestRequestDecompressionMiddleware_BufferedPath_SmallGzip
--- PASS: TestRequestDecompressionMiddleware_BufferedPath_SmallGzip (0.00s)
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_LargeGzip
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_LargeGzip (0.03s)
=== RUN   TestTracingMiddleware_StreamingRootSpanEndsAfterLLMSpan
--- PASS: TestTracingMiddleware_StreamingRootSpanEndsAfterLLMSpan (0.01s)
=== RUN   TestCollectDimensionHeaders
--- PASS: TestCollectDimensionHeaders (0.00s)
=== RUN   TestTracingMiddleware_SetsCorrelationHeaders
=== RUN   TestTracingMiddleware_SetsCorrelationHeaders/generates_request_id_when_absent
--- PASS: TestTracingMiddleware_SetsCorrelationHeaders/generates_request_id_when_absent (0.00s)
=== RUN   TestTracingMiddleware_SetsCorrelationHeaders/echoes_caller-supplied_request_id
--- PASS: TestTracingMiddleware_SetsCorrelationHeaders/echoes_caller-supplied_request_id (0.00s)
=== RUN   TestTracingMiddleware_SetsCorrelationHeaders/headers_survive_the_error_path
--- PASS: TestTracingMiddleware_SetsCorrelationHeaders/headers_survive_the_error_path (0.00s)
--- PASS: TestTracingMiddleware_SetsCorrelationHeaders (0.00s)
=== RUN   TestTracingMiddleware_AccessLogIncludesRequestID
--- PASS: TestTracingMiddleware_AccessLogIncludesRequestID (0.00s)
=== RUN   TestTransportPreAuthInterceptorMiddleware_HeaderVisibleToNext
--- PASS: TestTransportPreAuthInterceptorMiddleware_HeaderVisibleToNext (0.00s)
=== RUN   TestTransportPreAuthInterceptorMiddleware_BodyIsModifiable
--- PASS: TestTransportPreAuthInterceptorMiddleware_BodyIsModifiable (0.00s)
=== RUN   TestTransportPreAuthInterceptorMiddleware_ShortCircuitResponse
--- PASS: TestTransportPreAuthInterceptorMiddleware_ShortCircuitResponse (0.00s)
=== RUN   TestTransportPreAuthInterceptorMiddleware_PluginError
--- PASS: TestTransportPreAuthInterceptorMiddleware_PluginError (0.00s)
=== RUN   TestTransportPreAuthInterceptorMiddleware_PathMutationRejected
--- PASS: TestTransportPreAuthInterceptorMiddleware_PathMutationRejected (0.00s)
=== RUN   TestTransportPreAuthInterceptorMiddleware_NoPlugins
--- PASS: TestTransportPreAuthInterceptorMiddleware_NoPlugins (0.00s)
=== RUN   TestSecurityHeadersMiddleware_APINoStore
=== RUN   TestSecurityHeadersMiddleware_APINoStore/api_path_gets_no-store
--- PASS: TestSecurityHeadersMiddleware_APINoStore/api_path_gets_no-store (0.00s)
=== RUN   TestSecurityHeadersMiddleware_APINoStore/api_path_keeps_handler_policy
--- PASS: TestSecurityHeadersMiddleware_APINoStore/api_path_keeps_handler_policy (0.00s)
=== RUN   TestSecurityHeadersMiddleware_APINoStore/non-api_path_untouched
--- PASS: TestSecurityHeadersMiddleware_APINoStore/non-api_path_untouched (0.00s)
=== RUN   TestSecurityHeadersMiddleware_APINoStore/non-api_path_keeps_handler_policy
--- PASS: TestSecurityHeadersMiddleware_APINoStore/non-api_path_keeps_handler_policy (0.00s)
=== RUN   TestSecurityHeadersMiddleware_APINoStore/prefix_must_match_a_segment
--- PASS: TestSecurityHeadersMiddleware_APINoStore/prefix_must_match_a_segment (0.00s)
--- PASS: TestSecurityHeadersMiddleware_APINoStore (0.00s)
=== RUN   TestGetModelParameters_ResolvesQualifiedAndBareIDs
    model_parameters_test.go:28: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/model_parameters_test.go:28
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestGetModelParameters_ResolvesQualifiedAndBareIDs
--- FAIL: TestGetModelParameters_ResolvesQualifiedAndBareIDs (0.00s)
=== RUN   TestNotificationListFiltersByRole
--- PASS: TestNotificationListFiltersByRole (0.01s)
=== RUN   TestNotificationServicePublishValidatesAndPersists
--- PASS: TestNotificationServicePublishValidatesAndPersists (0.00s)
=== RUN   TestNotificationListStopsAtScanBudget
--- PASS: TestNotificationListStopsAtScanBudget (0.00s)
=== RUN   TestNotificationCreateRequiresLocalAdmin
--- PASS: TestNotificationCreateRequiresLocalAdmin (0.00s)
=== RUN   TestNotificationVisibleToRole
--- PASS: TestNotificationVisibleToRole (0.00s)
=== RUN   TestNotificationCursorRoundTrip
--- PASS: TestNotificationCursorRoundTrip (0.00s)
=== RUN   TestCreatePlugin_RejectsCustomPathWhenAuthBypassed
--- PASS: TestCreatePlugin_RejectsCustomPathWhenAuthBypassed (0.00s)
=== RUN   TestCreatePlugin_AllowsCustomPathWhenNotBypassed
--- PASS: TestCreatePlugin_AllowsCustomPathWhenNotBypassed (0.00s)
=== RUN   TestUpdatePlugin_RejectsCustomPathWhenAuthBypassed
--- PASS: TestUpdatePlugin_RejectsCustomPathWhenAuthBypassed (0.00s)
=== RUN   TestRestoreRedacted_OTELProfilesHeaders
--- PASS: TestRestoreRedacted_OTELProfilesHeaders (0.00s)
=== RUN   TestRestoreRedacted_KafkaSecretVarObjects
--- PASS: TestRestoreRedacted_KafkaSecretVarObjects (0.00s)
=== RUN   TestRestoreRedacted_FullyRedactedSentinel
--- PASS: TestRestoreRedacted_FullyRedactedSentinel (0.00s)
=== RUN   TestUpdatePlugin_ConfigMerge
--- PASS: TestUpdatePlugin_ConfigMerge (0.00s)
=== RUN   TestUpdatePlugin_ConfigMerge_NewPlugin
--- PASS: TestUpdatePlugin_ConfigMerge_NewPlugin (0.00s)
=== RUN   TestGetLoadedPlugins
--- PASS: TestGetLoadedPlugins (0.00s)
=== RUN   TestUpdatePricingOverride_ReplacesFullBody
    pricing_override_test.go:105: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:105
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestUpdatePricingOverride_ReplacesFullBody
--- FAIL: TestUpdatePricingOverride_ReplacesFullBody (0.00s)
=== RUN   TestPromptMutationsReloadCache
=== RUN   TestPromptMutationsReloadCache/deleteFolder
--- PASS: TestPromptMutationsReloadCache/deleteFolder (0.00s)
=== RUN   TestPromptMutationsReloadCache/createPrompt
--- PASS: TestPromptMutationsReloadCache/createPrompt (0.00s)
=== RUN   TestPromptMutationsReloadCache/updatePrompt
--- PASS: TestPromptMutationsReloadCache/updatePrompt (0.00s)
=== RUN   TestPromptMutationsReloadCache/deletePrompt
--- PASS: TestPromptMutationsReloadCache/deletePrompt (0.00s)
=== RUN   TestPromptMutationsReloadCache/createVersion
--- PASS: TestPromptMutationsReloadCache/createVersion (0.00s)
=== RUN   TestPromptMutationsReloadCache/deleteVersion
--- PASS: TestPromptMutationsReloadCache/deleteVersion (0.00s)
=== RUN   TestPromptMutationsReloadCache/createSession
--- PASS: TestPromptMutationsReloadCache/createSession (0.00s)
=== RUN   TestPromptMutationsReloadCache/updateSession
--- PASS: TestPromptMutationsReloadCache/updateSession (0.00s)
=== RUN   TestPromptMutationsReloadCache/deleteSession
--- PASS: TestPromptMutationsReloadCache/deleteSession (0.00s)
=== RUN   TestPromptMutationsReloadCache/renameSession
--- PASS: TestPromptMutationsReloadCache/renameSession (0.00s)
=== RUN   TestPromptMutationsReloadCache/commitSession
--- PASS: TestPromptMutationsReloadCache/commitSession (0.00s)
--- PASS: TestPromptMutationsReloadCache (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/deleteFolder
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/deleteFolder (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/createPrompt
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/createPrompt (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/updatePrompt
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/updatePrompt (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/deletePrompt
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/deletePrompt (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/createVersion
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/createVersion (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/deleteVersion
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/deleteVersion (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/createSession
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/createSession (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/updateSession
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/updateSession (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/deleteSession
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/deleteSession (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/renameSession
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/renameSession (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/commitSession
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/commitSession (0.00s)
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure (0.00s)
=== RUN   TestPromptMutationsSurviveReloaderProblems
=== RUN   TestPromptMutationsSurviveReloaderProblems/nil_reloader
--- PASS: TestPromptMutationsSurviveReloaderProblems/nil_reloader (0.00s)
=== RUN   TestPromptMutationsSurviveReloaderProblems/failing_reloader
--- PASS: TestPromptMutationsSurviveReloaderProblems/failing_reloader (0.00s)
--- PASS: TestPromptMutationsSurviveReloaderProblems (0.00s)
=== RUN   TestDatabricksKey_RedactedUpdateRoundTrip
    provider_keys_databricks_test.go:36: redacted workspace_url="https://dbc-1.cloud.databricks.com" api_format="ai_gateway" tags=true client_id="dbx-************************alue"
--- PASS: TestDatabricksKey_RedactedUpdateRoundTrip (0.00s)
=== RUN   TestDatabricksKey_Validate
=== RUN   TestDatabricksKey_Validate/pat_ok
--- PASS: TestDatabricksKey_Validate/pat_ok (0.00s)
=== RUN   TestDatabricksKey_Validate/oauth_ok
--- PASS: TestDatabricksKey_Validate/oauth_ok (0.00s)
=== RUN   TestDatabricksKey_Validate/nil_config
--- PASS: TestDatabricksKey_Validate/nil_config (0.00s)
=== RUN   TestDatabricksKey_Validate/no_workspace_url
--- PASS: TestDatabricksKey_Validate/no_workspace_url (0.00s)
=== RUN   TestDatabricksKey_Validate/no_auth_at_all
--- PASS: TestDatabricksKey_Validate/no_auth_at_all (0.00s)
=== RUN   TestDatabricksKey_Validate/half_a_service_principal
--- PASS: TestDatabricksKey_Validate/half_a_service_principal (0.00s)
=== RUN   TestDatabricksKey_Validate/half_a_pair_alongside_a_pat
--- PASS: TestDatabricksKey_Validate/half_a_pair_alongside_a_pat (0.00s)
=== RUN   TestDatabricksKey_Validate/bad_api_format
--- PASS: TestDatabricksKey_Validate/bad_api_format (0.00s)
=== RUN   TestDatabricksKey_Validate/env-ref_workspace_url
--- PASS: TestDatabricksKey_Validate/env-ref_workspace_url (0.00s)
=== RUN   TestDatabricksKey_Validate/env-ref_pat
--- PASS: TestDatabricksKey_Validate/env-ref_pat (0.00s)
--- PASS: TestDatabricksKey_Validate (0.00s)
=== RUN   TestDatabricksKey_SwitchToPAT
--- PASS: TestDatabricksKey_SwitchToPAT (0.00s)
=== RUN   TestMergeUpdatedKey_Value
=== RUN   TestMergeUpdatedKey_Value/echoed_current_redaction_preserves_stored_value
--- PASS: TestMergeUpdatedKey_Value/echoed_current_redaction_preserves_stored_value (0.00s)
=== RUN   TestMergeUpdatedKey_Value/mismatched_mask_still_preserves_stored_value
--- PASS: TestMergeUpdatedKey_Value/mismatched_mask_still_preserves_stored_value (0.00s)
=== RUN   TestMergeUpdatedKey_Value/genuine_new_plaintext_value_is_applied
--- PASS: TestMergeUpdatedKey_Value/genuine_new_plaintext_value_is_applied (0.00s)
=== RUN   TestMergeUpdatedKey_Value/genuine_env_ref_is_applied_not_preserved
--- PASS: TestMergeUpdatedKey_Value/genuine_env_ref_is_applied_not_preserved (0.00s)
=== RUN   TestMergeUpdatedKey_Value/empty_value_is_not_treated_as_redacted
--- PASS: TestMergeUpdatedKey_Value/empty_value_is_not_treated_as_redacted (0.00s)
--- PASS: TestMergeUpdatedKey_Value (0.00s)
=== RUN   TestMergeUpdatedKey_Name
=== RUN   TestMergeUpdatedKey_Name/name_omitted_from_update_preserves_stored_name
--- PASS: TestMergeUpdatedKey_Name/name_omitted_from_update_preserves_stored_name (0.00s)
=== RUN   TestMergeUpdatedKey_Name/explicit_new_name_is_applied
--- PASS: TestMergeUpdatedKey_Name/explicit_new_name_is_applied (0.00s)
--- PASS: TestMergeUpdatedKey_Name (0.00s)
=== RUN   TestMergeUpdatedKey_ProviderConfigMaskedPreviews
--- PASS: TestMergeUpdatedKey_ProviderConfigMaskedPreviews (0.00s)
=== RUN   TestMergeUpdatedKey_GithubCopilotKeyConfig
=== RUN   TestMergeUpdatedKey_GithubCopilotKeyConfig/masks_resolve_back_to_the_stored_values
--- PASS: TestMergeUpdatedKey_GithubCopilotKeyConfig/masks_resolve_back_to_the_stored_values (0.04s)
=== RUN   TestMergeUpdatedKey_GithubCopilotKeyConfig/a_mask_with_no_stored_counterpart_is_rejected
--- PASS: TestMergeUpdatedKey_GithubCopilotKeyConfig/a_mask_with_no_stored_counterpart_is_rejected (0.00s)
=== RUN   TestMergeUpdatedKey_GithubCopilotKeyConfig/unmasked_literals_overwrite_the_stored_values
--- PASS: TestMergeUpdatedKey_GithubCopilotKeyConfig/unmasked_literals_overwrite_the_stored_values (0.00s)
--- PASS: TestMergeUpdatedKey_GithubCopilotKeyConfig (0.04s)
=== RUN   TestMergeUpdatedKey_RejectsMaskWithoutStoredCounterpart
=== RUN   TestMergeUpdatedKey_RejectsMaskWithoutStoredCounterpart/missing_config_section
--- PASS: TestMergeUpdatedKey_RejectsMaskWithoutStoredCounterpart/missing_config_section (0.00s)
=== RUN   TestMergeUpdatedKey_RejectsMaskWithoutStoredCounterpart/missing_optional_field
--- PASS: TestMergeUpdatedKey_RejectsMaskWithoutStoredCounterpart/missing_optional_field (0.00s)
=== RUN   TestMergeUpdatedKey_RejectsMaskWithoutStoredCounterpart/empty_stored_value
--- PASS: TestMergeUpdatedKey_RejectsMaskWithoutStoredCounterpart/empty_stored_value (0.00s)
--- PASS: TestMergeUpdatedKey_RejectsMaskWithoutStoredCounterpart (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields
=== RUN   TestValidateProviderKeyRequiredNestedFields/azure_missing_endpoint
--- PASS: TestValidateProviderKeyRequiredNestedFields/azure_missing_endpoint (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/azure_nil_config
--- PASS: TestValidateProviderKeyRequiredNestedFields/azure_nil_config (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/azure_ok
--- PASS: TestValidateProviderKeyRequiredNestedFields/azure_ok (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/bedrock_missing_region
--- PASS: TestValidateProviderKeyRequiredNestedFields/bedrock_missing_region (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/bedrock_ok
--- PASS: TestValidateProviderKeyRequiredNestedFields/bedrock_ok (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/mantle_missing_region
--- PASS: TestValidateProviderKeyRequiredNestedFields/mantle_missing_region (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/mantle_ok
--- PASS: TestValidateProviderKeyRequiredNestedFields/mantle_ok (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/vllm_missing_url
--- PASS: TestValidateProviderKeyRequiredNestedFields/vllm_missing_url (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/vllm_missing_model_name
--- PASS: TestValidateProviderKeyRequiredNestedFields/vllm_missing_model_name (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/vllm_ok
--- PASS: TestValidateProviderKeyRequiredNestedFields/vllm_ok (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/openai_unaffected
--- PASS: TestValidateProviderKeyRequiredNestedFields/openai_unaffected (0.00s)
--- PASS: TestValidateProviderKeyRequiredNestedFields (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats
=== RUN   TestValidateProviderKeyGithubCopilotFormats/valid_literal_app_config
--- PASS: TestValidateProviderKeyGithubCopilotFormats/valid_literal_app_config (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/direct_token_needs_no_app_config
--- PASS: TestValidateProviderKeyGithubCopilotFormats/direct_token_needs_no_app_config (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/non-numeric_installation_id
--- PASS: TestValidateProviderKeyGithubCopilotFormats/non-numeric_installation_id (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/path_traversal_in_installation_id
--- PASS: TestValidateProviderKeyGithubCopilotFormats/path_traversal_in_installation_id (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/non-numeric_repository_id
--- PASS: TestValidateProviderKeyGithubCopilotFormats/non-numeric_repository_id (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/private_key_that_is_not_PEM
--- PASS: TestValidateProviderKeyGithubCopilotFormats/private_key_that_is_not_PEM (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/valid_PKCS#8_private_key
--- PASS: TestValidateProviderKeyGithubCopilotFormats/valid_PKCS#8_private_key (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/EC_private_key
--- PASS: TestValidateProviderKeyGithubCopilotFormats/EC_private_key (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/encrypted_private_key_envelope
--- PASS: TestValidateProviderKeyGithubCopilotFormats/encrypted_private_key_envelope (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/well-formed_envelope_with_a_garbage_DER_payload
--- PASS: TestValidateProviderKeyGithubCopilotFormats/well-formed_envelope_with_a_garbage_DER_payload (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/PEM_whose_newlines_survived_as_literal_backslash-n
--- PASS: TestValidateProviderKeyGithubCopilotFormats/PEM_whose_newlines_survived_as_literal_backslash-n (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/no_token_and_no_app_config_at_all
--- PASS: TestValidateProviderKeyGithubCopilotFormats/no_token_and_no_app_config_at_all (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/no_token,_missing_app_id
--- PASS: TestValidateProviderKeyGithubCopilotFormats/no_token,_missing_app_id (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/no_token,_missing_installation_id
--- PASS: TestValidateProviderKeyGithubCopilotFormats/no_token,_missing_installation_id (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/no_token,_missing_repository_id
--- PASS: TestValidateProviderKeyGithubCopilotFormats/no_token,_missing_repository_id (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/no_token,_missing_private_key
--- PASS: TestValidateProviderKeyGithubCopilotFormats/no_token,_missing_private_key (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/direct_token_alongside_an_incomplete_app_config
--- PASS: TestValidateProviderKeyGithubCopilotFormats/direct_token_alongside_an_incomplete_app_config (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/direct_token_alongside_a_malformed_app_config
--- PASS: TestValidateProviderKeyGithubCopilotFormats/direct_token_alongside_a_malformed_app_config (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/direct_token_alongside_a_complete_app_config
--- PASS: TestValidateProviderKeyGithubCopilotFormats/direct_token_alongside_a_complete_app_config (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/client_ID_as_app_id_is_accepted
--- PASS: TestValidateProviderKeyGithubCopilotFormats/client_ID_as_app_id_is_accepted (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/env_references_are_not_format-checked
--- PASS: TestValidateProviderKeyGithubCopilotFormats/env_references_are_not_format-checked (0.00s)
--- PASS: TestValidateProviderKeyGithubCopilotFormats (0.00s)
=== RUN   TestRefreshProviderModels_DelegatesToModelsManager
--- PASS: TestRefreshProviderModels_DelegatesToModelsManager (0.00s)
=== RUN   TestRefreshProviderKeyModels_DelegatesToModelsManager
--- PASS: TestRefreshProviderKeyModels_DelegatesToModelsManager (0.00s)
=== RUN   TestRefreshProviderModels_InFlightReturns409
--- PASS: TestRefreshProviderModels_InFlightReturns409 (0.00s)
=== RUN   TestRefreshProviderKeyModels_UnknownKeyReturns404
--- PASS: TestRefreshProviderKeyModels_UnknownKeyReturns404 (0.00s)
=== RUN   TestRefreshProviderModels_UnknownProviderReturns404
--- PASS: TestRefreshProviderModels_UnknownProviderReturns404 (0.00s)
=== RUN   TestCreateProviderKey_CustomBedrockRequiresRegion
--- PASS: TestCreateProviderKey_CustomBedrockRequiresRegion (0.02s)
=== RUN   TestApplyProviderConfigUpdates_OmittedBlocksArePreserved
--- PASS: TestApplyProviderConfigUpdates_OmittedBlocksArePreserved (0.01s)
=== RUN   TestApplyProviderConfigUpdates_ExplicitNullClears
--- PASS: TestApplyProviderConfigUpdates_ExplicitNullClears (0.00s)
=== RUN   TestApplyProviderConfigUpdates_PresentBlocksAreReplaced
--- PASS: TestApplyProviderConfigUpdates_PresentBlocksAreReplaced (0.00s)
=== RUN   TestAddProvider_ReloadsRuntimeEvenWhenModelDiscoveryIsSkipped
--- PASS: TestAddProvider_ReloadsRuntimeEvenWhenModelDiscoveryIsSkipped (0.01s)
=== RUN   TestAddProvider_ReturnsErrorWhenRuntimeReloadFails
--- PASS: TestAddProvider_ReturnsErrorWhenRuntimeReloadFails (0.00s)
=== RUN   TestUpdateProvider_RejectsKeysInBody
--- PASS: TestUpdateProvider_RejectsKeysInBody (0.02s)
=== RUN   TestUpdateProvider_PassesThroughForEmptyOrAbsentKeys
=== RUN   TestUpdateProvider_PassesThroughForEmptyOrAbsentKeys/keys_field_omitted_entirely
--- PASS: TestUpdateProvider_PassesThroughForEmptyOrAbsentKeys/keys_field_omitted_entirely (0.00s)
=== RUN   TestUpdateProvider_PassesThroughForEmptyOrAbsentKeys/keys_explicitly_null
--- PASS: TestUpdateProvider_PassesThroughForEmptyOrAbsentKeys/keys_explicitly_null (0.00s)
=== RUN   TestUpdateProvider_PassesThroughForEmptyOrAbsentKeys/keys_explicitly_empty_array
--- PASS: TestUpdateProvider_PassesThroughForEmptyOrAbsentKeys/keys_explicitly_empty_array (0.00s)
--- PASS: TestUpdateProvider_PassesThroughForEmptyOrAbsentKeys (0.00s)
=== RUN   TestListModels_UnknownKeysDoNotFilter
--- PASS: TestListModels_UnknownKeysDoNotFilter (0.00s)
=== RUN   TestListModels_ReturnsExactAccessibleByKeysAndSkipsDisabledKeys
--- PASS: TestListModels_ReturnsExactAccessibleByKeysAndSkipsDisabledKeys (0.00s)
=== RUN   TestListModels_AppliesQueryAndLimitAfterFiltering
--- PASS: TestListModels_AppliesQueryAndLimitAfterFiltering (0.00s)
=== RUN   TestListModels_MarksDeprecatedModelsWithoutFiltering
    bedrock_test.go:95: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/integrations/bedrock_test.go:95
        	Error:      	Received unexpected error:
        	            	failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestGenericRouter_MarkDeprecatedListModelsResponseUsesCatalog808912714\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestGenericRouter_MarkDeprecatedListModelsResponseUsesCatalog808912714\\001\\pricing.json" after host
        	Test:       	TestGenericRouter_MarkDeprecatedListModelsResponseUsesCatalog
--- FAIL: TestGenericRouter_MarkDeprecatedListModelsResponseUsesCatalog (7.01s)
=== RUN   Test_parseS3URI
=== RUN   Test_parseS3URI/full_S3_URI_with_key
--- PASS: Test_parseS3URI/full_S3_URI_with_key (0.00s)
=== RUN   Test_parseS3URI/S3_URI_with_bucket_only
--- PASS: Test_parseS3URI/S3_URI_with_bucket_only (0.00s)
=== RUN   Test_parseS3URI/S3_URI_with_bucket_no_trailing_slash
--- PASS: Test_parseS3URI/S3_URI_with_bucket_no_trailing_slash (0.00s)
=== RUN   Test_parseS3URI/plain_bucket_name
--- PASS: Test_parseS3URI/plain_bucket_name (0.00s)
=== RUN   Test_parseS3URI/S3_URI_with_nested_key
--- PASS: Test_parseS3URI/S3_URI_with_nested_key (0.00s)
=== RUN   Test_parseS3URI/empty_string
--- PASS: Test_parseS3URI/empty_string (0.00s)
--- PASS: Test_parseS3URI (0.00s)
=== RUN   Test_createBedrockRouteConfigs
--- PASS: Test_createBedrockRouteConfigs (0.00s)
=== RUN   Test_createBedrockConverseRouteConfig
--- PASS: Test_createBedrockConverseRouteConfig (0.00s)
=== RUN   Test_createBedrockConverseStreamRouteConfig
--- PASS: Test_createBedrockConverseStreamRouteConfig (0.00s)
=== RUN   Test_createBedrockInvokeRouteConfig
--- PASS: Test_createBedrockInvokeRouteConfig (0.00s)
=== RUN   TestBedrockInvokeEmbeddingResponseLangChainCohereAlias
--- PASS: TestBedrockInvokeEmbeddingResponseLangChainCohereAlias (0.00s)
=== RUN   TestBedrockInvokeTypedCohereResponseLangChainAlias
--- PASS: TestBedrockInvokeTypedCohereResponseLangChainAlias (0.00s)
=== RUN   TestBedrockInvokeTypedCohereResponseWithoutFloatIsNotAliased
--- PASS: TestBedrockInvokeTypedCohereResponseWithoutFloatIsNotAliased (0.00s)
=== RUN   TestBedrockInvokeEmbeddingResponseLangChainLeavesTitanUnchanged
--- PASS: TestBedrockInvokeEmbeddingResponseLangChainLeavesTitanUnchanged (0.00s)
=== RUN   TestBedrockInvokeTypedTitanResponseLangChainLeavesEnvelopeUnchanged
--- PASS: TestBedrockInvokeTypedTitanResponseLangChainLeavesEnvelopeUnchanged (0.00s)
=== RUN   TestBedrockInvokeTypedCohereResponseUsesResolvedModelForAlias
--- PASS: TestBedrockInvokeTypedCohereResponseUsesResolvedModelForAlias (0.00s)
=== RUN   Test_createBedrockInvokeWithResponseStreamRouteConfig
--- PASS: Test_createBedrockInvokeWithResponseStreamRouteConfig (0.00s)
=== RUN   Test_bedrockStreamErrorConverterEncodesEventStreamException
--- PASS: Test_bedrockStreamErrorConverterEncodesEventStreamException (0.00s)
=== RUN   Test_toBedrockEventStreamExceptionAcceptsBedrockError
--- PASS: Test_toBedrockEventStreamExceptionAcceptsBedrockError (0.00s)
=== RUN   Test_handleStreamingBedrockUnknownErrorResponseFallsBackToEventStreamException
--- PASS: Test_handleStreamingBedrockUnknownErrorResponseFallsBackToEventStreamException (0.00s)
=== RUN   Test_createBedrockRerankRouteConfig
--- PASS: Test_createBedrockRerankRouteConfig (0.00s)
=== RUN   Test_createBedrockRerankResponseConverterEmitsBedrockShape
--- PASS: Test_createBedrockRerankResponseConverterEmitsBedrockShape (0.00s)
=== RUN   Test_createBedrockRerankResponseConverterUsesRawResponse
--- PASS: Test_createBedrockRerankResponseConverterUsesRawResponse (0.00s)
=== RUN   Test_createBedrockRerankResponseConverterConvertsForCrossProvider
--- PASS: Test_createBedrockRerankResponseConverterConvertsForCrossProvider (0.00s)
=== RUN   Test_createBedrockRerankRouteRequestConverter
--- PASS: Test_createBedrockRerankRouteRequestConverter (0.00s)
=== RUN   Test_createBedrockRouteConfigsIncludesRerankForCompositePrefixes
--- PASS: Test_createBedrockRouteConfigsIncludesRerankForCompositePrefixes (0.00s)
=== RUN   Test_createBedrockBatchRouteConfigs
--- PASS: Test_createBedrockBatchRouteConfigs (0.00s)
=== RUN   Test_createBedrockFilesRouteConfigs
--- PASS: Test_createBedrockFilesRouteConfigs (0.00s)
=== RUN   Test_parseS3PutObjectRequest
=== RUN   Test_parseS3PutObjectRequest/valid_request
--- PASS: Test_parseS3PutObjectRequest/valid_request (0.00s)
=== RUN   Test_parseS3PutObjectRequest/simple_key_without_folder
--- PASS: Test_parseS3PutObjectRequest/simple_key_without_folder (0.00s)
=== RUN   Test_parseS3PutObjectRequest/missing_bucket
--- PASS: Test_parseS3PutObjectRequest/missing_bucket (0.00s)
=== RUN   Test_parseS3PutObjectRequest/missing_key
--- PASS: Test_parseS3PutObjectRequest/missing_key (0.00s)
--- PASS: Test_parseS3PutObjectRequest (0.00s)
=== RUN   Test_parseS3PutObjectRequest_invalidType
--- PASS: Test_parseS3PutObjectRequest_invalidType (0.00s)
=== RUN   Test_s3PutObjectPostCallback
=== RUN   Test_s3PutObjectPostCallback/valid_response_with_ID
--- PASS: Test_s3PutObjectPostCallback/valid_response_with_ID (0.00s)
=== RUN   Test_s3PutObjectPostCallback/nil_response
--- PASS: Test_s3PutObjectPostCallback/nil_response (0.00s)
--- PASS: Test_s3PutObjectPostCallback (0.00s)
=== RUN   Test_s3GetObjectPostCallback
=== RUN   Test_s3GetObjectPostCallback/valid_response
--- PASS: Test_s3GetObjectPostCallback/valid_response (0.00s)
=== RUN   Test_s3GetObjectPostCallback/nil_response
--- PASS: Test_s3GetObjectPostCallback/nil_response (0.00s)
--- PASS: Test_s3GetObjectPostCallback (0.00s)
=== RUN   Test_s3HeadObjectPostCallback
=== RUN   Test_s3HeadObjectPostCallback/valid_response
--- PASS: Test_s3HeadObjectPostCallback/valid_response (0.00s)
=== RUN   Test_s3HeadObjectPostCallback/nil_response
--- PASS: Test_s3HeadObjectPostCallback/nil_response (0.00s)
--- PASS: Test_s3HeadObjectPostCallback (0.00s)
=== RUN   Test_s3DeleteObjectPostCallback
--- PASS: Test_s3DeleteObjectPostCallback (0.00s)
=== RUN   Test_s3ListObjectsV2PostCallback
--- PASS: Test_s3ListObjectsV2PostCallback (0.00s)
=== RUN   Test_extractBedrockBatchListQueryParams
=== RUN   Test_extractBedrockBatchListQueryParams/all_params
--- PASS: Test_extractBedrockBatchListQueryParams/all_params (0.00s)
=== RUN   Test_extractBedrockBatchListQueryParams/no_params
--- PASS: Test_extractBedrockBatchListQueryParams/no_params (0.00s)
=== RUN   Test_extractBedrockBatchListQueryParams/invalid_maxResults
--- PASS: Test_extractBedrockBatchListQueryParams/invalid_maxResults (0.00s)
--- PASS: Test_extractBedrockBatchListQueryParams (0.00s)
=== RUN   Test_extractBedrockJobArnFromPath
=== RUN   Test_extractBedrockJobArnFromPath/valid_job_ARN_for_Bedrock
--- PASS: Test_extractBedrockJobArnFromPath/valid_job_ARN_for_Bedrock (0.00s)
=== RUN   Test_extractBedrockJobArnFromPath/URL_encoded_job_ARN
--- PASS: Test_extractBedrockJobArnFromPath/URL_encoded_job_ARN (0.00s)
=== RUN   Test_extractBedrockJobArnFromPath/non-Bedrock_provider_strips_ARN_prefix
--- PASS: Test_extractBedrockJobArnFromPath/non-Bedrock_provider_strips_ARN_prefix (0.00s)
=== RUN   Test_extractBedrockJobArnFromPath/missing_job_arn
--- PASS: Test_extractBedrockJobArnFromPath/missing_job_arn (0.00s)
=== RUN   Test_extractBedrockJobArnFromPath/empty_job_arn
--- PASS: Test_extractBedrockJobArnFromPath/empty_job_arn (0.00s)
--- PASS: Test_extractBedrockJobArnFromPath (0.00s)
=== RUN   Test_extractS3ListObjectsV2Params
=== RUN   Test_extractS3ListObjectsV2Params/all_params
--- PASS: Test_extractS3ListObjectsV2Params/all_params (0.00s)
=== RUN   Test_extractS3ListObjectsV2Params/bucket_only
--- PASS: Test_extractS3ListObjectsV2Params/bucket_only (0.00s)
=== RUN   Test_extractS3ListObjectsV2Params/missing_bucket
--- PASS: Test_extractS3ListObjectsV2Params/missing_bucket (0.00s)
--- PASS: Test_extractS3ListObjectsV2Params (0.00s)
=== RUN   Test_extractS3BucketKeyFromPath
=== RUN   Test_extractS3BucketKeyFromPath/content_operation
--- PASS: Test_extractS3BucketKeyFromPath/content_operation (0.00s)
=== RUN   Test_extractS3BucketKeyFromPath/missing_bucket
--- PASS: Test_extractS3BucketKeyFromPath/missing_bucket (0.00s)
=== RUN   Test_extractS3BucketKeyFromPath/missing_key
--- PASS: Test_extractS3BucketKeyFromPath/missing_key (0.00s)
--- PASS: Test_extractS3BucketKeyFromPath (0.00s)
=== RUN   TestCreateCohereRouteConfigsIncludesRerank
--- PASS: TestCreateCohereRouteConfigsIncludesRerank (0.00s)
=== RUN   TestCohereRerankRouteRequestConverter
--- PASS: TestCohereRerankRouteRequestConverter (0.00s)
=== RUN   TestCohereRerankResponseConverterEmitsCohereShape
--- PASS: TestCohereRerankResponseConverterEmitsCohereShape (0.00s)
=== RUN   TestCohereRerankResponseConverterUsesRawResponse
--- PASS: TestCohereRerankResponseConverterUsesRawResponse (0.00s)
=== RUN   TestCohereRerankResponseConverterConvertsForCrossProvider
--- PASS: TestCohereRerankResponseConverterConvertsForCrossProvider (0.00s)
=== RUN   TestCohereChatResponseConverterEmitsCohereV2Shape
--- PASS: TestCohereChatResponseConverterEmitsCohereV2Shape (0.01s)
=== RUN   TestCohereChatResponseConverterUsesNativeRawResponse
--- PASS: TestCohereChatResponseConverterUsesNativeRawResponse (0.00s)
=== RUN   TestCohereChatResponseConverterConvertsCrossProviderRawResponse
--- PASS: TestCohereChatResponseConverterConvertsCrossProviderRawResponse (0.00s)
=== RUN   TestExtractModelAndRequestType_LargePayloadUsesMetadataWithoutBodyParse
--- PASS: TestExtractModelAndRequestType_LargePayloadUsesMetadataWithoutBodyParse (0.00s)
=== RUN   TestExtractModelAndRequestType_LargeBodyHeuristicSkipsParse
--- PASS: TestExtractModelAndRequestType_LargeBodyHeuristicSkipsParse (0.00s)
=== RUN   TestRewriteGenAIRawRequestBodyRedactsOnlyContentFields
--- PASS: TestRewriteGenAIRawRequestBodyRedactsOnlyContentFields (0.00s)
=== RUN   TestRewriteGenAIRawRequestBodySupportsCountTokensEnvelope
--- PASS: TestRewriteGenAIRawRequestBodySupportsCountTokensEnvelope (0.00s)
=== RUN   TestRewriteGenAIRawRequestBodyRejectsUnmappedLiteral
--- PASS: TestRewriteGenAIRawRequestBodyRejectsUnmappedLiteral (0.00s)
=== RUN   TestRewriteGenAIRawRequestBodyRejectsSensitiveObjectKeys
--- PASS: TestRewriteGenAIRawRequestBodyRejectsSensitiveObjectKeys (0.00s)
=== RUN   TestCreateGenAIRerankRouteConfig
--- PASS: TestCreateGenAIRerankRouteConfig (0.00s)
=== RUN   TestCreateGenAIRouteConfigsIncludesRerank
--- PASS: TestCreateGenAIRouteConfigsIncludesRerank (0.00s)
=== RUN   TestGenAISpeechStreamResponseConverter
--- PASS: TestGenAISpeechStreamResponseConverter (0.00s)
=== RUN   TestGenAISpeechStreamDoneResponseIncludesUsageAndFinishReason
--- PASS: TestGenAISpeechStreamDoneResponseIncludesUsageAndFinishReason (0.00s)
=== RUN   TestGenAICamelCaseSpeechVoiceSurvivesConversion
--- PASS: TestGenAICamelCaseSpeechVoiceSurvivesConversion (0.09s)
=== RUN   TestExtractAndSetModelAndRequestTypePreservesRawBodyForGenerateContent
--- PASS: TestExtractAndSetModelAndRequestTypePreservesRawBodyForGenerateContent (0.00s)
=== RUN   TestExtractAndSetModelAndRequestTypeNoRawPassthroughWithoutExplicitGemini
--- PASS: TestExtractAndSetModelAndRequestTypeNoRawPassthroughWithoutExplicitGemini (0.00s)
=== RUN   TestExtractAndSetModelAndRequestTypeDoesNotRawPassthroughEmbedding
--- PASS: TestExtractAndSetModelAndRequestTypeDoesNotRawPassthroughEmbedding (0.01s)
=== RUN   TestExtractAndSetModelAndRequestTypeDetectsImageEditAtAnyPartPosition
=== RUN   TestExtractAndSetModelAndRequestTypeDetectsImageEditAtAnyPartPosition/image_first
--- PASS: TestExtractAndSetModelAndRequestTypeDetectsImageEditAtAnyPartPosition/image_first (0.00s)
=== RUN   TestExtractAndSetModelAndRequestTypeDetectsImageEditAtAnyPartPosition/text_first
--- PASS: TestExtractAndSetModelAndRequestTypeDetectsImageEditAtAnyPartPosition/text_first (0.00s)
--- PASS: TestExtractAndSetModelAndRequestTypeDetectsImageEditAtAnyPartPosition (0.00s)
=== RUN   TestGenAIBatchCreateConverterCarriesRawBody
--- PASS: TestGenAIBatchCreateConverterCarriesRawBody (0.05s)
=== RUN   TestGenAIBatchCreateConverterCarriesDisplayNameAndInlineTools
--- PASS: TestGenAIBatchCreateConverterCarriesDisplayNameAndInlineTools (0.06s)
=== RUN   TestGenAICachedContentCreateParserRejectsNonStringScalars
--- PASS: TestGenAICachedContentCreateParserRejectsNonStringScalars (0.00s)
=== RUN   TestGenAICachedContentCreateParserCarriesRawBody
--- PASS: TestGenAICachedContentCreateParserCarriesRawBody (0.00s)
=== RUN   TestCreateGenAIRouteConfigsIncludesRerankForCompositePrefixes
--- PASS: TestCreateGenAIRouteConfigsIncludesRerankForCompositePrefixes (0.00s)
=== RUN   TestGenAIRerankRequestConverter
--- PASS: TestGenAIRerankRequestConverter (0.00s)
=== RUN   TestGenAIRerankResponseConverterRestoresCallerRecordIDs
--- PASS: TestGenAIRerankResponseConverterRestoresCallerRecordIDs (0.00s)
=== RUN   TestGenAIRerankRequestConverterRequestsDocuments
--- PASS: TestGenAIRerankRequestConverterRequestsDocuments (0.00s)
=== RUN   TestCreateGenAIRouteConfigsIncludesModelMetadataRoute
--- PASS: TestCreateGenAIRouteConfigsIncludesModelMetadataRoute (0.00s)
=== RUN   TestExtractGeminiModelMetadataParams
--- PASS: TestExtractGeminiModelMetadataParams (0.00s)
=== RUN   TestConvertGeminiModelMetadataResponse
--- PASS: TestConvertGeminiModelMetadataResponse (0.00s)
=== RUN   TestConvertGeminiModelMetadataResponse_MatchesRequestedModelNotFirst
--- PASS: TestConvertGeminiModelMetadataResponse_MatchesRequestedModelNotFirst (0.00s)
=== RUN   TestConvertGeminiModelMetadataResponse_EmptyReturnsMinimalModel
--- PASS: TestConvertGeminiModelMetadataResponse_EmptyReturnsMinimalModel (0.00s)
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting/gemini_to_gemini/generateContent
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting/gemini_to_gemini/generateContent (0.00s)
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting/gemini_to_gemini/streamGenerateContent
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting/gemini_to_gemini/streamGenerateContent (0.00s)
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting/gemini_to_vertex/generateContent
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting/gemini_to_vertex/generateContent (0.00s)
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting/gemini_to_vertex/streamGenerateContent
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting/gemini_to_vertex/streamGenerateContent (0.00s)
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting/vertex_to_gemini/generateContent
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting/vertex_to_gemini/generateContent (0.00s)
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting/vertex_to_gemini/streamGenerateContent
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting/vertex_to_gemini/streamGenerateContent (0.00s)
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting/vertex_to_vertex/generateContent
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting/vertex_to_vertex/generateContent (0.00s)
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting/vertex_to_vertex/streamGenerateContent
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting/vertex_to_vertex/streamGenerateContent (0.00s)
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting (0.00s)
=== RUN   TestOpenAIWireCostResponse
=== RUN   TestOpenAIWireCostResponse/chat_renders_cost_as_float_total_without_mutating_the_shared_response
--- PASS: TestOpenAIWireCostResponse/chat_renders_cost_as_float_total_without_mutating_the_shared_response (0.01s)
=== RUN   TestOpenAIWireCostResponse/responses_renders_cost_as_float_total
--- PASS: TestOpenAIWireCostResponse/responses_renders_cost_as_float_total (0.02s)
=== RUN   TestOpenAIWireCostResponse/image_renders_cost_as_float_total
--- PASS: TestOpenAIWireCostResponse/image_renders_cost_as_float_total (0.00s)
=== RUN   TestOpenAIWireCostResponse/chat_without_cost_passes_through_the_same_pointer
--- PASS: TestOpenAIWireCostResponse/chat_without_cost_passes_through_the_same_pointer (0.00s)
=== RUN   TestOpenAIWireCostResponse/raw_upstream_payload_passes_through_unchanged
--- PASS: TestOpenAIWireCostResponse/raw_upstream_payload_passes_through_unchanged (0.00s)
--- PASS: TestOpenAIWireCostResponse (0.03s)
=== RUN   TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter
=== RUN   TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/both_fields_present
--- PASS: TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/both_fields_present (0.00s)
=== RUN   TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/neither_field_present
--- PASS: TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/neither_field_present (0.00s)
=== RUN   TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/anchor_only
--- PASS: TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/anchor_only (0.00s)
=== RUN   TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/seconds_only
--- PASS: TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/seconds_only (0.00s)
=== RUN   TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/out-of-range_seconds
--- PASS: TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/out-of-range_seconds (0.00s)
=== RUN   TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/non-numeric_seconds
--- PASS: TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/non-numeric_seconds (0.00s)
--- PASS: TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter (0.00s)
=== RUN   TestParseTranscriptionMultipartRequest_ExtraParamsPassthrough
--- PASS: TestParseTranscriptionMultipartRequest_ExtraParamsPassthrough (0.00s)
=== RUN   TestParseTranscriptionMultipartRequest_TypedFieldsNotShadowedByExtraParams
--- PASS: TestParseTranscriptionMultipartRequest_TypedFieldsNotShadowedByExtraParams (0.00s)
=== RUN   TestParseTranscriptionMultipartRequest_PreservesFilename
--- PASS: TestParseTranscriptionMultipartRequest_PreservesFilename (0.00s)
=== RUN   TestForwardPassthroughRequestHeaderKeepsCookie
--- PASS: TestForwardPassthroughRequestHeaderKeepsCookie (0.00s)
=== RUN   Test_handleStreamingSSESendsHeartbeatDuringIdleGap
--- PASS: Test_handleStreamingSSESendsHeartbeatDuringIdleGap (0.00s)
=== RUN   Test_handleStreamingGenAIEmitsNoHeartbeat
--- PASS: Test_handleStreamingGenAIEmitsNoHeartbeat (0.00s)
=== RUN   Test_passthroughHeartbeatEligible
=== RUN   Test_passthroughHeartbeatEligible/plain_SSE
--- PASS: Test_passthroughHeartbeatEligible/plain_SSE (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/SSE_with_charset_param
--- PASS: Test_passthroughHeartbeatEligible/SSE_with_charset_param (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/SSE_uppercase
--- PASS: Test_passthroughHeartbeatEligible/SSE_uppercase (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/SSE_with_surrounding_space_before_param
--- PASS: Test_passthroughHeartbeatEligible/SSE_with_surrounding_space_before_param (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/raw_JSON_passthrough
--- PASS: Test_passthroughHeartbeatEligible/raw_JSON_passthrough (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/empty_content-type
--- PASS: Test_passthroughHeartbeatEligible/empty_content-type (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/binary_passthrough
--- PASS: Test_passthroughHeartbeatEligible/binary_passthrough (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/prefix-collision_suffix
--- PASS: Test_passthroughHeartbeatEligible/prefix-collision_suffix (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/prefix-collision_longer_type
--- PASS: Test_passthroughHeartbeatEligible/prefix-collision_longer_type (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/Gemini_SSE_passthrough
--- PASS: Test_passthroughHeartbeatEligible/Gemini_SSE_passthrough (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/Vertex_SSE_passthrough
--- PASS: Test_passthroughHeartbeatEligible/Vertex_SSE_passthrough (0.00s)
--- PASS: Test_passthroughHeartbeatEligible (0.00s)
=== RUN   TestCreateHandler_SkipsRequestParserInLargePayloadMode
--- PASS: TestCreateHandler_SkipsRequestParserInLargePayloadMode (0.00s)
=== RUN   TestCreateHandler_UsesRequestParserWhenNotInLargePayloadMode
--- PASS: TestCreateHandler_UsesRequestParserWhenNotInLargePayloadMode (0.00s)
=== RUN   TestResolveLargePayloadMetadata_NilContext
--- PASS: TestResolveLargePayloadMetadata_NilContext (0.00s)
=== RUN   TestResolveLargePayloadMetadata_SyncPath
--- PASS: TestResolveLargePayloadMetadata_SyncPath (0.00s)
=== RUN   TestResolveLargePayloadMetadata_DeferredReady
--- PASS: TestResolveLargePayloadMetadata_DeferredReady (0.00s)
=== RUN   TestResolveLargePayloadMetadata_DeferredNotReady
--- PASS: TestResolveLargePayloadMetadata_DeferredNotReady (0.00s)
=== RUN   TestResolveLargePayloadMetadata_SyncTakesPrecedence
--- PASS: TestResolveLargePayloadMetadata_SyncTakesPrecedence (0.00s)
=== RUN   TestListModelsNarrowsTheFanOut
--- PASS: TestListModelsNarrowsTheFanOut (0.00s)
=== RUN   TestListModelsLeavesFanOutAloneWhenNothingResolved
--- PASS: TestListModelsLeavesFanOutAloneWhenNothingResolved (0.00s)
=== RUN   TestListModelsWithoutAccessResolver
--- PASS: TestListModelsWithoutAccessResolver (0.00s)
=== RUN   TestGrantedProvidersShapeMatchesWhatTheRouterPublishes
--- PASS: TestGrantedProvidersShapeMatchesWhatTheRouterPublishes (0.00s)
=== RUN   Test_handleStreamingInterceptionErrorUsesErrorConverter
--- PASS: Test_handleStreamingInterceptionErrorUsesErrorConverter (0.00s)
=== RUN   Test_handleStreamingInterceptionErrorAnthropicSSEFraming
--- PASS: Test_handleStreamingInterceptionErrorAnthropicSSEFraming (0.00s)
=== RUN   Test_handleStreamingInterceptionErrorIsSanitized
--- PASS: Test_handleStreamingInterceptionErrorIsSanitized (0.00s)
=== RUN   Test_handleStreamingInterceptionErrorBedrockEventStream
--- PASS: Test_handleStreamingInterceptionErrorBedrockEventStream (0.00s)
=== RUN   Test_handleStreamingPlainInterceptionErrorKeepsFlatFormat
--- PASS: Test_handleStreamingPlainInterceptionErrorKeepsFlatFormat (0.00s)
=== RUN   Test_handleStreamingMissingSpeechConverterReturnsError
--- PASS: Test_handleStreamingMissingSpeechConverterReturnsError (0.00s)
=== RUN   TestIsPassthroughRequestOnlyMatchesRegisteredNativeRoutes
--- PASS: TestIsPassthroughRequestOnlyMatchesRegisteredNativeRoutes (0.00s)
=== RUN   TestParsePassthroughBody_MultipartExtractsModelAfterFilePart
--- PASS: TestParsePassthroughBody_MultipartExtractsModelAfterFilePart (0.00s)
=== RUN   TestChatGPTPassthroughRouterRegistersCodexResponsesPost
--- PASS: TestChatGPTPassthroughRouterRegistersCodexResponsesPost (0.00s)
=== RUN   TestChatGPTBackgroundRoutes
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/codex/models?client_version=0.159.3
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/codex/models?client_version=0.159.3 (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/HEAD/backend-api/codex/models
--- PASS: TestChatGPTBackgroundRoutes/HEAD/backend-api/codex/models (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/ps/plugins/list
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/ps/plugins/list (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/POST/backend-api/ps/plugins/install
--- PASS: TestChatGPTBackgroundRoutes/POST/backend-api/ps/plugins/install (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/ps/mcp
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/ps/mcp (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/PUT/backend-api/ps/mcp/config
--- PASS: TestChatGPTBackgroundRoutes/PUT/backend-api/ps/mcp/config (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/wham/usage
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/wham/usage (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/plugins/featured
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/plugins/featured (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/POST/backend-api/codex/analytics-events/events
--- PASS: TestChatGPTBackgroundRoutes/POST/backend-api/codex/analytics-events/events (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/codex/analytics-events/events
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/codex/analytics-events/events (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/POST/backend-api/codex/models
--- PASS: TestChatGPTBackgroundRoutes/POST/backend-api/codex/models (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/DELETE/backend-api/accounts
--- PASS: TestChatGPTBackgroundRoutes/DELETE/backend-api/accounts (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/ps/mcp-lookalike
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/ps/mcp-lookalike (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/ps/plugins-other/list
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/ps/plugins-other/list (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/TRACE/backend-api/wham/usage
--- PASS: TestChatGPTBackgroundRoutes/TRACE/backend-api/wham/usage (0.00s)
--- PASS: TestChatGPTBackgroundRoutes (0.00s)
=== RUN   TestChatGPTBackgroundWireAndHooks
--- PASS: TestChatGPTBackgroundWireAndHooks (0.01s)
=== RUN   TestChatGPTUpstreamURLDefaultsToChatGPT
--- PASS: TestChatGPTUpstreamURLDefaultsToChatGPT (0.00s)
=== RUN   TestChatGPTUpstreamURLHonoursDiagnosticOverride
--- PASS: TestChatGPTUpstreamURLHonoursDiagnosticOverride (0.00s)
=== RUN   TestRunwarePassthroughRouterRegistersCatchAll
--- PASS: TestRunwarePassthroughRouterRegistersCatchAll (0.00s)
=== RUN   TestRequestWithSettableExtraParams_OpenAIChatRequest
=== RUN   TestRequestWithSettableExtraParams_OpenAIChatRequest/SetExtraParams_populates_both_standalone_and_embedded_ExtraParams
--- PASS: TestRequestWithSettableExtraParams_OpenAIChatRequest/SetExtraParams_populates_both_standalone_and_embedded_ExtraParams (0.00s)
=== RUN   TestRequestWithSettableExtraParams_OpenAIChatRequest/extra_params_propagate_through_ToBifrostChatRequest
--- PASS: TestRequestWithSettableExtraParams_OpenAIChatRequest/extra_params_propagate_through_ToBifrostChatRequest (0.00s)
--- PASS: TestRequestWithSettableExtraParams_OpenAIChatRequest (0.00s)
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIChatRequest_implements_RequestWithSettableExtraParams
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIChatRequest_implements_RequestWithSettableExtraParams (0.00s)
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAITextCompletionRequest_implements_RequestWithSettableExtraParams
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAITextCompletionRequest_implements_RequestWithSettableExtraParams (0.00s)
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIResponsesRequest_implements_RequestWithSettableExtraParams
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIResponsesRequest_implements_RequestWithSettableExtraParams (0.00s)
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIEmbeddingRequest_implements_RequestWithSettableExtraParams
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIEmbeddingRequest_implements_RequestWithSettableExtraParams (0.00s)
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAISpeechRequest_implements_RequestWithSettableExtraParams
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAISpeechRequest_implements_RequestWithSettableExtraParams (0.00s)
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIImageGenerationRequest_implements_RequestWithSettableExtraParams
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIImageGenerationRequest_implements_RequestWithSettableExtraParams (0.00s)
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIImageEditRequest_implements_RequestWithSettableExtraParams
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIImageEditRequest_implements_RequestWithSettableExtraParams (0.00s)
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIImageVariationRequest_implements_RequestWithSettableExtraParams
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIImageVariationRequest_implements_RequestWithSettableExtraParams (0.00s)
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes (0.01s)
=== RUN   TestExtraParamsRequiresPassthroughHeader
=== RUN   TestExtraParamsRequiresPassthroughHeader/extra_params_NOT_extracted_without_passthrough_header
--- PASS: TestExtraParamsRequiresPassthroughHeader/extra_params_NOT_extracted_without_passthrough_header (0.07s)
=== RUN   TestExtraParamsRequiresPassthroughHeader/extra_params_extracted_with_passthrough_header
--- PASS: TestExtraParamsRequiresPassthroughHeader/extra_params_extracted_with_passthrough_header (0.00s)
--- PASS: TestExtraParamsRequiresPassthroughHeader (0.07s)
=== RUN   TestExtraParamsPassthrough_NestedStructures
--- PASS: TestExtraParamsPassthrough_NestedStructures (0.00s)
=== RUN   TestExtraParamsPassthrough_EndToEnd
--- PASS: TestExtraParamsPassthrough_EndToEnd (0.00s)
=== RUN   TestExtraParamsPassthrough_NoExtraParamsKey
--- PASS: TestExtraParamsPassthrough_NoExtraParamsKey (0.00s)
=== RUN   TestOpenAIChatStructuredOutputRequestParserAndConverter
--- PASS: TestOpenAIChatStructuredOutputRequestParserAndConverter (0.00s)
=== RUN   TestCreateHandler_AnthropicRouteSetsPassthroughFlags
--- PASS: TestCreateHandler_AnthropicRouteSetsPassthroughFlags (0.00s)
=== RUN   TestCreateHandler_AnthropicThreadContinueRefused
--- PASS: TestCreateHandler_AnthropicThreadContinueRefused (0.00s)
=== RUN   TestCreateHandler_AnthropicThreadContinueRefusedStreaming
--- PASS: TestCreateHandler_AnthropicThreadContinueRefusedStreaming (0.00s)
=== RUN   TestCreateHandler_AnthropicThreadCreateProceeds
--- PASS: TestCreateHandler_AnthropicThreadCreateProceeds (0.00s)
=== RUN   TestCreateHandler_CustomParserFailureClosesConnection
--- PASS: TestCreateHandler_CustomParserFailureClosesConnection (0.00s)
=== RUN   TestCreateHandler_DefaultJSONParserFailureClosesConnection
--- PASS: TestCreateHandler_DefaultJSONParserFailureClosesConnection (0.00s)
=== RUN   TestCreateHandler_ParseFailureClosesKeepAliveSocket
=== RUN   TestCreateHandler_ParseFailureClosesKeepAliveSocket/valid_keep-alive_requests_reuse_the_socket
--- PASS: TestCreateHandler_ParseFailureClosesKeepAliveSocket/valid_keep-alive_requests_reuse_the_socket (0.00s)
=== RUN   TestCreateHandler_ParseFailureClosesKeepAliveSocket/malformed_request_closes_the_socket
--- PASS: TestCreateHandler_ParseFailureClosesKeepAliveSocket/malformed_request_closes_the_socket (0.00s)
--- PASS: TestCreateHandler_ParseFailureClosesKeepAliveSocket (0.00s)
=== RUN   TestExtraParamsSetViaInterfaceMutatesOriginalReq
--- PASS: TestExtraParamsSetViaInterfaceMutatesOriginalReq (0.00s)
=== RUN   TestExtractModelFromPath
=== RUN   TestExtractModelFromPath/azure_deployment_chat
--- PASS: TestExtractModelFromPath/azure_deployment_chat (0.00s)
=== RUN   TestExtractModelFromPath/azure_deployment_leading-stripped
--- PASS: TestExtractModelFromPath/azure_deployment_leading-stripped (0.00s)
=== RUN   TestExtractModelFromPath/genai_models_with_action
--- PASS: TestExtractModelFromPath/genai_models_with_action (0.00s)
=== RUN   TestExtractModelFromPath/genai_models_stream_action
--- PASS: TestExtractModelFromPath/genai_models_stream_action (0.00s)
=== RUN   TestExtractModelFromPath/genai_tunedModels
--- PASS: TestExtractModelFromPath/genai_tunedModels (0.00s)
=== RUN   TestExtractModelFromPath/vertex_fully-qualified
--- PASS: TestExtractModelFromPath/vertex_fully-qualified (0.00s)
=== RUN   TestExtractModelFromPath/no_model_segment
--- PASS: TestExtractModelFromPath/no_model_segment (0.00s)
=== RUN   TestExtractModelFromPath/deployments_with_no_trailing_segment
--- PASS: TestExtractModelFromPath/deployments_with_no_trailing_segment (0.00s)
--- PASS: TestExtractModelFromPath (0.00s)
=== RUN   TestExtractPassthroughModel
=== RUN   TestExtractPassthroughModel/azure_deployment_path_overrides_empty_body
--- PASS: TestExtractPassthroughModel/azure_deployment_path_overrides_empty_body (0.00s)
=== RUN   TestExtractPassthroughModel/body_fallback_when_path_has_no_model
--- PASS: TestExtractPassthroughModel/body_fallback_when_path_has_no_model (0.00s)
=== RUN   TestExtractPassthroughModel/path_wins_over_body
--- PASS: TestExtractPassthroughModel/path_wins_over_body (0.00s)
=== RUN   TestExtractPassthroughModel/both_empty
--- PASS: TestExtractPassthroughModel/both_empty (0.00s)
=== RUN   TestExtractPassthroughModel/vertex_resource_body_model
--- PASS: TestExtractPassthroughModel/vertex_resource_body_model (0.00s)
=== RUN   TestExtractPassthroughModel/genai_resource_body_model
--- PASS: TestExtractPassthroughModel/genai_resource_body_model (0.00s)
=== RUN   TestExtractPassthroughModel/slashed_non-resource_model_untouched
--- PASS: TestExtractPassthroughModel/slashed_non-resource_model_untouched (0.00s)
--- PASS: TestExtractPassthroughModel (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_AnthropicOAuthForwarded
--- PASS: TestApplyPassthroughCallerAuth_AnthropicOAuthForwarded (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_OpenAIJWTForwarded
--- PASS: TestApplyPassthroughCallerAuth_OpenAIJWTForwarded (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_LoopbackHTTPOverrideForwardsJWT
--- PASS: TestApplyPassthroughCallerAuth_LoopbackHTTPOverrideForwardsJWT (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_APIKeysStayStripped
=== RUN   TestApplyPassthroughCallerAuth_APIKeysStayStripped/openai_plain_api_key
--- PASS: TestApplyPassthroughCallerAuth_APIKeysStayStripped/openai_plain_api_key (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_APIKeysStayStripped/anthropic_api_key_bearer
--- PASS: TestApplyPassthroughCallerAuth_APIKeysStayStripped/anthropic_api_key_bearer (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_APIKeysStayStripped/provider_override_bedrock
--- PASS: TestApplyPassthroughCallerAuth_APIKeysStayStripped/provider_override_bedrock (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_APIKeysStayStripped/openai_two-segment_token
--- PASS: TestApplyPassthroughCallerAuth_APIKeysStayStripped/openai_two-segment_token (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_APIKeysStayStripped/http_upstream_override
--- PASS: TestApplyPassthroughCallerAuth_APIKeysStayStripped/http_upstream_override (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_APIKeysStayStripped/http_lookalike_host
--- PASS: TestApplyPassthroughCallerAuth_APIKeysStayStripped/http_lookalike_host (0.00s)
--- PASS: TestApplyPassthroughCallerAuth_APIKeysStayStripped (0.00s)
=== RUN   TestExtractAndParseFallbacks_GeminiGenerationRequest
--- PASS: TestExtractAndParseFallbacks_GeminiGenerationRequest (0.00s)
=== RUN   TestJSONErrorsPreserveBodyForBodylessStatus
--- PASS: TestJSONErrorsPreserveBodyForBodylessStatus (0.00s)
=== RUN   TestSendStreamError_PropagatesProviderStatusCode
=== RUN   TestSendStreamError_PropagatesProviderStatusCode/provider_400_-_Bedrock_ValidationException_/_OpenAI_invalid_request_error
--- PASS: TestSendStreamError_PropagatesProviderStatusCode/provider_400_-_Bedrock_ValidationException_/_OpenAI_invalid_request_error (0.00s)
=== RUN   TestSendStreamError_PropagatesProviderStatusCode/provider_429_-_rate_limiting_(all_providers)
--- PASS: TestSendStreamError_PropagatesProviderStatusCode/provider_429_-_rate_limiting_(all_providers) (0.00s)
=== RUN   TestSendStreamError_PropagatesProviderStatusCode/provider_503_-_Bedrock_ServiceUnavailableException
--- PASS: TestSendStreamError_PropagatesProviderStatusCode/provider_503_-_Bedrock_ServiceUnavailableException (0.00s)
=== RUN   TestSendStreamError_PropagatesProviderStatusCode/provider_529_-_Anthropic_overloaded_error
--- PASS: TestSendStreamError_PropagatesProviderStatusCode/provider_529_-_Anthropic_overloaded_error (0.00s)
=== RUN   TestSendStreamError_PropagatesProviderStatusCode/nil_StatusCode_defaults_to_500
--- PASS: TestSendStreamError_PropagatesProviderStatusCode/nil_StatusCode_defaults_to_500 (0.00s)
--- PASS: TestSendStreamError_PropagatesProviderStatusCode (0.00s)
=== RUN   TestSendStreamError_OpenAIErrorFormat
--- PASS: TestSendStreamError_OpenAIErrorFormat (0.00s)
=== RUN   TestSendStreamError_AnthropicErrorFormat
--- PASS: TestSendStreamError_AnthropicErrorFormat (0.00s)
=== RUN   TestSendStreamError_BedrockErrorFormat
--- PASS: TestSendStreamError_BedrockErrorFormat (0.00s)
=== RUN   TestSendStreamError_ForwardsProviderHeaders
--- PASS: TestSendStreamError_ForwardsProviderHeaders (0.00s)
=== RUN   TestTryStreamLargeResponse_AppliesRoutedIdentityHeaders
--- PASS: TestTryStreamLargeResponse_AppliesRoutedIdentityHeaders (0.00s)
FAIL
FAIL	github.com/maximhq/bifrost/transports/bifrost-http/integrations	7.677s
    providers_test.go:549: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModels_MarksDeprecatedModelsWithoutFiltering1455123442\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModels_MarksDeprecatedModelsWithoutFiltering1455123442\\001\\pricing.json" after host
--- FAIL: TestListModels_MarksDeprecatedModelsWithoutFiltering (7.00s)
=== RUN   TestListBaseModels_IncludesDeprecatedPricingRows
    providers_test.go:592: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListBaseModels_IncludesDeprecatedPricingRows707984320\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListBaseModels_IncludesDeprecatedPricingRows707984320\\001\\pricing.json" after host
--- FAIL: TestListBaseModels_IncludesDeprecatedPricingRows (7.00s)
=== RUN   TestEnrichListModelsResponse_MarksDeprecatedPricingRows
    providers_test.go:617: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestEnrichListModelsResponse_MarksDeprecatedPricingRows1794733899\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestEnrichListModelsResponse_MarksDeprecatedPricingRows1794733899\\001\\pricing.json" after host
--- FAIL: TestEnrichListModelsResponse_MarksDeprecatedPricingRows (7.00s)
=== RUN   TestListModels_UnfilteredIgnoresKeys
--- PASS: TestListModels_UnfilteredIgnoresKeys (0.00s)
=== RUN   TestListModels_UnfilteredWithoutKeysReturnsAllUnfilteredModels
--- PASS: TestListModels_UnfilteredWithoutKeysReturnsAllUnfilteredModels (0.00s)
=== RUN   TestListModelDetails_ErrorsWhenModelCatalogUnavailable
--- PASS: TestListModelDetails_ErrorsWhenModelCatalogUnavailable (0.00s)
=== RUN   TestListModelDetails_UnknownKeysDoNotFilter
--- PASS: TestListModelDetails_UnknownKeysDoNotFilter (0.00s)
=== RUN   TestListModelDetails_SkipsUnknownKeysAndFiltersWithValid
--- PASS: TestListModelDetails_SkipsUnknownKeysAndFiltersWithValid (0.00s)
=== RUN   TestListModelDetails_SkipsDisabledKeysAndFiltersWithValid
--- PASS: TestListModelDetails_SkipsDisabledKeysAndFiltersWithValid (0.00s)
=== RUN   TestListModelDetails_UnfilteredIgnoresKeys
--- PASS: TestListModelDetails_UnfilteredIgnoresKeys (0.00s)
=== RUN   TestListModelDetails_IncludesPricing
    providers_test.go:888: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_IncludesPricing2208179822\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_IncludesPricing2208179822\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_IncludesPricing (7.00s)
=== RUN   TestListModelDetails_ResolvesCatalogPricing
=== RUN   TestListModelDetails_ResolvesCatalogPricing/Together_catalog_provider
    providers_test.go:1027: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingTogether_catalog_pro150855590\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingTogether_catalog_pro150855590\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_ResolvesCatalogPricing/Together_catalog_provider (7.01s)
=== RUN   TestListModelDetails_ResolvesCatalogPricing/Together_alias
    providers_test.go:1027: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingTogether_alias4267945468\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingTogether_alias4267945468\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_ResolvesCatalogPricing/Together_alias (7.02s)
=== RUN   TestListModelDetails_ResolvesCatalogPricing/Azure_alias
    providers_test.go:1027: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingAzure_alias184254088\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingAzure_alias184254088\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_ResolvesCatalogPricing/Azure_alias (7.02s)
=== RUN   TestListModelDetails_ResolvesCatalogPricing/Azure_alias_with_empty_model_name
    providers_test.go:1027: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingAzure_alias_with_emp1144308229\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingAzure_alias_with_emp1144308229\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_ResolvesCatalogPricing/Azure_alias_with_empty_model_name (7.00s)
=== RUN   TestListModelDetails_ResolvesCatalogPricing/Azure_alias_falls_back_to_alias_key
    providers_test.go:1027: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingAzure_alias_falls_ba3807605529\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingAzure_alias_falls_ba3807605529\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_ResolvesCatalogPricing/Azure_alias_falls_back_to_alias_key (7.00s)
--- FAIL: TestListModelDetails_ResolvesCatalogPricing (35.06s)
=== RUN   TestListModelDetails_AppliesGlobalOverrideWithoutMutatingBase
    providers_test.go:1094: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_AppliesGlobalOverrideWithoutMutatingBase1434450568\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_AppliesGlobalOverrideWithoutMutatingBase1434450568\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_AppliesGlobalOverrideWithoutMutatingBase (7.00s)
=== RUN   TestListModelDetails_OverrideIndexIsDeduplicated
    providers_test.go:1150: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_OverrideIndexIsDeduplicated1271822688\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_OverrideIndexIsDeduplicated1271822688\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_OverrideIndexIsDeduplicated (7.00s)
=== RUN   TestListModelDetails_NoOverridesOmitsNewFields
    providers_test.go:1182: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_NoOverridesOmitsNewFields1624610472\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_NoOverridesOmitsNewFields1624610472\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_NoOverridesOmitsNewFields (7.01s)
=== RUN   TestListModelDetails_PatchToSameValueIsNotMarkedOverridden
    providers_test.go:1197: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_PatchToSameValueIsNotMarkedOverridden3477404324\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_PatchToSameValueIsNotMarkedOverridden3477404324\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_PatchToSameValueIsNotMarkedOverridden (7.00s)
=== RUN   TestListModelDetails_OverrideWithoutBaseCatalogRow
    providers_test.go:1233: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_OverrideWithoutBaseCatalogRow1735290225\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_OverrideWithoutBaseCatalogRow1735290225\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_OverrideWithoutBaseCatalogRow (7.00s)
=== RUN   TestListModelDetails_VirtualKeyScopedOverrideIsInformationalOnly
    providers_test.go:1262: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_VirtualKeyScopedOverrideIsInformationalOnly2341351988\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_VirtualKeyScopedOverrideIsInformationalOnly2341351988\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_VirtualKeyScopedOverrideIsInformationalOnly (7.00s)
=== RUN   TestListModelDetails_AppliesProviderScopedOverride
    providers_test.go:1297: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_AppliesProviderScopedOverride2941577636\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_AppliesProviderScopedOverride2941577636\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_AppliesProviderScopedOverride (7.00s)
=== RUN   TestParseVKValueFromRequest
=== RUN   TestParseVKValueFromRequest/x-bf-vk_header
--- PASS: TestParseVKValueFromRequest/x-bf-vk_header (0.00s)
=== RUN   TestParseVKValueFromRequest/Authorization_Bearer_header
--- PASS: TestParseVKValueFromRequest/Authorization_Bearer_header (0.00s)
=== RUN   TestParseVKValueFromRequest/x-api-key_header
--- PASS: TestParseVKValueFromRequest/x-api-key_header (0.00s)
=== RUN   TestParseVKValueFromRequest/x-goog-api-key_header
--- PASS: TestParseVKValueFromRequest/x-goog-api-key_header (0.00s)
=== RUN   TestParseVKValueFromRequest/no_header_returns_empty_string
--- PASS: TestParseVKValueFromRequest/no_header_returns_empty_string (0.00s)
=== RUN   TestParseVKValueFromRequest/non-VK_Bearer_token_returns_empty_string
--- PASS: TestParseVKValueFromRequest/non-VK_Bearer_token_returns_empty_string (0.00s)
=== RUN   TestParseVKValueFromRequest/x-bf-vk_takes_priority_over_Authorization
--- PASS: TestParseVKValueFromRequest/x-bf-vk_takes_priority_over_Authorization (0.00s)
--- PASS: TestParseVKValueFromRequest (0.00s)
=== RUN   TestListModels_VKFilterHidesBlacklistedModel
--- PASS: TestListModels_VKFilterHidesBlacklistedModel (0.00s)
=== RUN   TestListModels_VKFilterUnionsDuplicateProviderConfigs
--- PASS: TestListModels_VKFilterUnionsDuplicateProviderConfigs (0.00s)
=== RUN   TestListModels_VKFilterRestrictsToAllowedProviderAndModels
--- PASS: TestListModels_VKFilterRestrictsToAllowedProviderAndModels (0.00s)
=== RUN   TestListModels_VKFilterAllowsAllModelsWithWildcard
--- PASS: TestListModels_VKFilterAllowsAllModelsWithWildcard (0.00s)
=== RUN   TestListModels_VKFilterDeniesAllModelsWhenAllowedModelsEmpty
--- PASS: TestListModels_VKFilterDeniesAllModelsWhenAllowedModelsEmpty (0.00s)
=== RUN   TestListModels_VKFilterNoProviderConfigsDeniesAll
--- PASS: TestListModels_VKFilterNoProviderConfigsDeniesAll (0.00s)
=== RUN   TestListModels_VKFilterBlockedExplicitProviderReturnsEmptyResult
--- PASS: TestListModels_VKFilterBlockedExplicitProviderReturnsEmptyResult (0.00s)
=== RUN   TestParseModelListQuery_VKAppliesResolvedAccess
--- PASS: TestParseModelListQuery_VKAppliesResolvedAccess (0.00s)
=== RUN   TestParseModelListQuery_LeavesListingUnfilteredWhenNothingResolved
=== RUN   TestParseModelListQuery_LeavesListingUnfilteredWhenNothingResolved/no_key_presented
--- PASS: TestParseModelListQuery_LeavesListingUnfilteredWhenNothingResolved/no_key_presented (0.00s)
=== RUN   TestParseModelListQuery_LeavesListingUnfilteredWhenNothingResolved/key_resolves_to_nothing
--- PASS: TestParseModelListQuery_LeavesListingUnfilteredWhenNothingResolved/key_resolves_to_nothing (0.00s)
--- PASS: TestParseModelListQuery_LeavesListingUnfilteredWhenNothingResolved (0.00s)
=== RUN   TestListModels_NoVKFilterReturnsAll
--- PASS: TestListModels_NoVKFilterReturnsAll (0.00s)
=== RUN   TestListModels_UsesCatalogAwareAliasMatchingForKeyAllowlist
--- PASS: TestListModels_UsesCatalogAwareAliasMatchingForKeyAllowlist (0.00s)
=== RUN   TestListModels_KeyModelAllowlistIsCaseInsensitive
--- PASS: TestListModels_KeyModelAllowlistIsCaseInsensitive (0.00s)
=== RUN   TestListModels_KeyBlacklistIsCaseInsensitive
--- PASS: TestListModels_KeyBlacklistIsCaseInsensitive (0.00s)
=== RUN   TestRealtimeSessionRoutesOnlyExposeGAClientSecrets
--- PASS: TestRealtimeSessionRoutesOnlyExposeGAClientSecrets (0.00s)
=== RUN   TestResolveRealtimeClientSecretTarget
=== PAUSE TestResolveRealtimeClientSecretTarget
=== RUN   TestResolveRealtimeClientSecretTarget_NormalizesModel
=== PAUSE TestResolveRealtimeClientSecretTarget_NormalizesModel
=== RUN   TestGATranscriptionSessionEndToEndThroughFullNormalizationPath
=== PAUSE TestGATranscriptionSessionEndToEndThroughFullNormalizationPath
=== RUN   TestFullRealtimeSessionWithTranscriptionSiblingEndToEnd
=== PAUSE TestFullRealtimeSessionWithTranscriptionSiblingEndToEnd
=== RUN   TestRewriteGASessionTranscriptionModelAppliesAliasResolution
=== PAUSE TestRewriteGASessionTranscriptionModelAppliesAliasResolution
=== RUN   TestRewriteGASessionTranscriptionModelNoOpForFullRealtimeSession
=== PAUSE TestRewriteGASessionTranscriptionModelNoOpForFullRealtimeSession
=== RUN   TestParseRealtimeEphemeralKeyMapping
=== PAUSE TestParseRealtimeEphemeralKeyMapping
=== RUN   TestParseRealtimeEphemeralKeyMapping_NestedFallback
=== PAUSE TestParseRealtimeEphemeralKeyMapping_NestedFallback
=== RUN   TestParseRealtimeEphemeralKeyMapping_RejectsCompositeShapes
=== PAUSE TestParseRealtimeEphemeralKeyMapping_RejectsCompositeShapes
=== RUN   TestRealtimeMappingVirtualKeyUsesRequestHeader
=== PAUSE TestRealtimeMappingVirtualKeyUsesRequestHeader
=== RUN   TestRealtimeMappingVirtualKeyUsesSettledBearerCredential
=== PAUSE TestRealtimeMappingVirtualKeyUsesSettledBearerCredential
=== RUN   TestRealtimeMappingVirtualKeyFallsBackToContextValue
=== PAUSE TestRealtimeMappingVirtualKeyFallsBackToContextValue
=== RUN   TestReplaceAndCacheRealtimeEphemeralToken
=== PAUSE TestReplaceAndCacheRealtimeEphemeralToken
=== RUN   TestReplaceAndCacheRealtimeEphemeralTokenRejectsExpiredSecret
=== PAUSE TestReplaceAndCacheRealtimeEphemeralTokenRejectsExpiredSecret
=== RUN   TestIsJSONContentType
=== PAUSE TestIsJSONContentType
=== RUN   TestRealtimeClientSecretsEvaluateMintingGovernance_RequiresAccess
=== PAUSE TestRealtimeClientSecretsEvaluateMintingGovernance_RequiresAccess
=== RUN   TestRealtimeClientSecretsEvaluateMintingGovernance_PassesContext
=== PAUSE TestRealtimeClientSecretsEvaluateMintingGovernance_PassesContext
=== RUN   TestRealtimeClientSecretsEvaluateMintingGovernance_ContinuesWithoutGovernance
=== PAUSE TestRealtimeClientSecretsEvaluateMintingGovernance_ContinuesWithoutGovernance
=== RUN   TestRealtimeTurnFinalEventUsesTranscriptionCompletionOnlyForTranscriptionSessions
=== PAUSE TestRealtimeTurnFinalEventUsesTranscriptionCompletionOnlyForTranscriptionSessions
=== RUN   TestRealtimeTurnCompletionContentModelsTranscriptAsOutputOnlyForTranscriptionSessions
=== PAUSE TestRealtimeTurnCompletionContentModelsTranscriptAsOutputOnlyForTranscriptionSessions
=== RUN   TestBuildRealtimeTurnPostResponseModelsTranscriptionAsOutput
=== PAUSE TestBuildRealtimeTurnPostResponseModelsTranscriptionAsOutput
=== RUN   TestBuildRealtimeTurnPreRequestIncludesResponseCreateContent
=== PAUSE TestBuildRealtimeTurnPreRequestIncludesResponseCreateContent
=== RUN   TestBuildRealtimeTurnPreRequestCombinesConversationAndResponseInput
=== PAUSE TestBuildRealtimeTurnPreRequestCombinesConversationAndResponseInput
=== RUN   TestBuildRealtimeTurnPreRequestPreservesEmptyResponseTools
=== PAUSE TestBuildRealtimeTurnPreRequestPreservesEmptyResponseTools
=== RUN   TestShouldAccumulateRealtimeOutput
--- PASS: TestShouldAccumulateRealtimeOutput (0.00s)
=== RUN   TestExtractRealtimeTurnSummary
--- PASS: TestExtractRealtimeTurnSummary (0.00s)
=== RUN   TestFinalizedRealtimeInputSummary
--- PASS: TestFinalizedRealtimeInputSummary (0.00s)
=== RUN   TestFinalizedRealtimeToolOutputSummary
--- PASS: TestFinalizedRealtimeToolOutputSummary (0.00s)
=== RUN   TestPendingRealtimeInputUpdate
=== PAUSE TestPendingRealtimeInputUpdate
=== RUN   TestPendingRealtimeToolOutputUpdate
=== PAUSE TestPendingRealtimeToolOutputUpdate
=== RUN   TestRealtimeSessionDedupeNestedRawEvents
=== PAUSE TestRealtimeSessionDedupeNestedRawEvents
=== RUN   TestBuildRealtimeTurnPostResponseUsesFullResponseDonePayload
--- PASS: TestBuildRealtimeTurnPostResponseUsesFullResponseDonePayload (0.01s)
=== RUN   TestBuildRealtimeTurnPostResponseMergesTextAndToolCalls
--- PASS: TestBuildRealtimeTurnPostResponseMergesTextAndToolCalls (0.00s)
=== RUN   TestFinalizeRealtimeTurnHooksWithErrorCompletesActiveHooks
=== PAUSE TestFinalizeRealtimeTurnHooksWithErrorCompletesActiveHooks
=== RUN   TestNewBifrostErrorFromRealtimeErrorCarriesRealtimeMetadata
=== PAUSE TestNewBifrostErrorFromRealtimeErrorCarriesRealtimeMetadata
=== RUN   TestRefuseUnauthenticatedRealtime_AnonymousRefusedWhenEnforced
--- PASS: TestRefuseUnauthenticatedRealtime_AnonymousRefusedWhenEnforced (0.00s)
=== RUN   TestRefuseUnauthenticatedRealtime_AnonymousAllowedWhenNotEnforced
--- PASS: TestRefuseUnauthenticatedRealtime_AnonymousAllowedWhenNotEnforced (0.00s)
=== RUN   TestRefuseUnauthenticatedRealtime_VirtualKeyAllowed
--- PASS: TestRefuseUnauthenticatedRealtime_VirtualKeyAllowed (0.00s)
=== RUN   TestRefuseUnauthenticatedRealtime_MappedEphemeralTokenAllowed
--- PASS: TestRefuseUnauthenticatedRealtime_MappedEphemeralTokenAllowed (0.00s)
=== RUN   TestRefuseUnauthenticatedRealtime_UnmappedBifrostPrefixedTokenAllowed
--- PASS: TestRefuseUnauthenticatedRealtime_UnmappedBifrostPrefixedTokenAllowed (0.00s)
=== RUN   TestRefuseUnauthenticatedRealtime_NonEphemeralBearerRefused
--- PASS: TestRefuseUnauthenticatedRealtime_NonEphemeralBearerRefused (0.00s)
=== RUN   TestRefuseUnauthenticatedRealtime_DirectKeyAllowed
--- PASS: TestRefuseUnauthenticatedRealtime_DirectKeyAllowed (0.00s)
=== RUN   TestWSRealtimeHandleUpgrade_AnonymousRefusedBeforeTargetResolution
--- PASS: TestWSRealtimeHandleUpgrade_AnonymousRefusedBeforeTargetResolution (0.00s)
=== RUN   TestWSRealtimeHandleUpgrade_TargetErrorStaysInBandWhenNotEnforced
--- PASS: TestWSRealtimeHandleUpgrade_TargetErrorStaysInBandWhenNotEnforced (0.00s)
=== RUN   TestRefuseUnresolvedRealtimeCredential_ForgedVirtualKeyRefused
--- PASS: TestRefuseUnresolvedRealtimeCredential_ForgedVirtualKeyRefused (0.00s)
=== RUN   TestRefuseUnresolvedRealtimeCredential_ResolvedVirtualKeyAdmitted
--- PASS: TestRefuseUnresolvedRealtimeCredential_ResolvedVirtualKeyAdmitted (0.00s)
=== RUN   TestRefuseUnresolvedRealtimeCredential_EphemeralTokenExempt
--- PASS: TestRefuseUnresolvedRealtimeCredential_EphemeralTokenExempt (0.00s)
=== RUN   TestRefuseUnresolvedRealtimeCredential_MappedEphemeralTokenNotExempt
--- PASS: TestRefuseUnresolvedRealtimeCredential_MappedEphemeralTokenNotExempt (0.00s)
=== RUN   TestRefuseUnresolvedRealtimeCredential_ProviderTokenOnlyMappingExempt
--- PASS: TestRefuseUnresolvedRealtimeCredential_ProviderTokenOnlyMappingExempt (0.00s)
=== RUN   TestRefuseUnresolvedRealtimeCredential_DirectKeyAdmitted
--- PASS: TestRefuseUnresolvedRealtimeCredential_DirectKeyAdmitted (0.00s)
=== RUN   TestRefuseUnresolvedRealtimeCredential_NotEnforcedAllowsEverything
--- PASS: TestRefuseUnresolvedRealtimeCredential_NotEnforcedAllowsEverything (0.00s)
=== RUN   TestSessionLoginInvalidPayloadDoesNotExposeDecoderDetails
--- PASS: TestSessionLoginInvalidPayloadDoesNotExposeDecoderDetails (0.00s)
=== RUN   TestPrepareRequestInvalidPayloadDoesNotExposeDecoderDetails
--- PASS: TestPrepareRequestInvalidPayloadDoesNotExposeDecoderDetails (0.00s)
=== RUN   TestComplexityAnalyzerConfigGetDegradesOnUnreadableConfig
    routing_test.go:80: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/routing_test.go:80
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestComplexityAnalyzerConfigGetDegradesOnUnreadableConfig
--- FAIL: TestComplexityAnalyzerConfigGetDegradesOnUnreadableConfig (0.00s)
=== RUN   TestComplexityAnalyzerConfigGetStillFailsWhenStoreUnreachable
--- PASS: TestComplexityAnalyzerConfigGetStillFailsWhenStoreUnreachable (0.00s)
=== RUN   TestComplexityAnalyzerConfigGetReturnsDefaultsWhenUnset
    routing_test.go:124: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/routing_test.go:124
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestComplexityAnalyzerConfigGetReturnsDefaultsWhenUnset
--- FAIL: TestComplexityAnalyzerConfigGetReturnsDefaultsWhenUnset (0.00s)
=== RUN   TestRetryComplexitySemanticWarmup
=== RUN   TestRetryComplexitySemanticWarmup/accepts_a_failed_warmup
    routing_test.go:154: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/routing_test.go:154
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestRetryComplexitySemanticWarmup/accepts_a_failed_warmup
--- FAIL: TestRetryComplexitySemanticWarmup/accepts_a_failed_warmup (0.00s)
=== RUN   TestRetryComplexitySemanticWarmup/rejects_a_non-failed_warmup
    routing_test.go:172: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/routing_test.go:172
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestRetryComplexitySemanticWarmup/rejects_a_non-failed_warmup
--- FAIL: TestRetryComplexitySemanticWarmup/rejects_a_non-failed_warmup (0.00s)
--- FAIL: TestRetryComplexitySemanticWarmup (0.00s)
=== RUN   TestComplexityAnalyzerConfigPutPersistsAndReloads
    routing_test.go:185: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/routing_test.go:185
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestComplexityAnalyzerConfigPutPersistsAndReloads
--- FAIL: TestComplexityAnalyzerConfigPutPersistsAndReloads (0.00s)
=== RUN   TestComplexityAnalyzerConfigPutRejectsInvalidPayloads
    routing_test.go:233: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/routing_test.go:233
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestComplexityAnalyzerConfigPutRejectsInvalidPayloads
--- FAIL: TestComplexityAnalyzerConfigPutRejectsInvalidPayloads (0.00s)
=== RUN   TestComplexityAnalyzerConfigResetPersistsDefaultsAndReloads
    routing_test.go:285: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/routing_test.go:285
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestComplexityAnalyzerConfigResetPersistsDefaultsAndReloads
--- FAIL: TestComplexityAnalyzerConfigResetPersistsDefaultsAndReloads (0.00s)
=== RUN   TestComplexityAnalyzerConfigResetReportsReloadFailure
    routing_test.go:393: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/routing_test.go:393
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestComplexityAnalyzerConfigResetReportsReloadFailure
--- FAIL: TestComplexityAnalyzerConfigResetReportsReloadFailure (0.00s)
=== RUN   TestRoutingRoutesServeCanonicalAndLegacyPaths
--- PASS: TestRoutingRoutesServeCanonicalAndLegacyPaths (0.00s)
=== RUN   TestCleanupOrphanSkillFilesDeletesDBFallbackBlobs
    skills_cleanup_test.go:48: new config store: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
--- FAIL: TestCleanupOrphanSkillFilesDeletesDBFallbackBlobs (0.00s)
=== RUN   TestCleanupOrphanSkillFilesDeletesOnlyUnreferencedUploadObjects
    skills_cleanup_test.go:102: new config store: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
--- FAIL: TestCleanupOrphanSkillFilesDeletesOnlyUnreferencedUploadObjects (0.00s)
=== RUN   TestSkillsServingGenericFileDownloadDecodesEncodedPathParams
    skills_serving_test.go:22: new config store: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
--- FAIL: TestSkillsServingGenericFileDownloadDecodesEncodedPathParams (0.00s)
=== RUN   TestClaudeMarketplaceGitRepoContainsMarketplaceAndCloneablePlugin
    skills_serving_test.go:88: new config store: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
--- FAIL: TestClaudeMarketplaceGitRepoContainsMarketplaceAndCloneablePlugin (0.00s)
=== RUN   TestPrepareSpeechRequestVoiceValidation
=== RUN   TestPrepareSpeechRequestVoiceValidation/elevenlabs_alias_without_voice_is_deferred_to_provider
--- PASS: TestPrepareSpeechRequestVoiceValidation/elevenlabs_alias_without_voice_is_deferred_to_provider (0.01s)
=== RUN   TestPrepareSpeechRequestVoiceValidation/elevenlabs_literal_sound_model_without_voice_is_allowed
--- PASS: TestPrepareSpeechRequestVoiceValidation/elevenlabs_literal_sound_model_without_voice_is_allowed (0.00s)
=== RUN   TestPrepareSpeechRequestVoiceValidation/elevenlabs_tts_without_voice_is_deferred_to_provider_(no_transport_400)
--- PASS: TestPrepareSpeechRequestVoiceValidation/elevenlabs_tts_without_voice_is_deferred_to_provider_(no_transport_400) (0.00s)
=== RUN   TestPrepareSpeechRequestVoiceValidation/non-elevenlabs_provider_without_voice_still_errors_at_the_transport
--- PASS: TestPrepareSpeechRequestVoiceValidation/non-elevenlabs_provider_without_voice_still_errors_at_the_transport (0.00s)
--- PASS: TestPrepareSpeechRequestVoiceValidation (0.01s)
=== RUN   TestSSEStreamReaderNoEventBatching
--- PASS: TestSSEStreamReaderNoEventBatching (0.01s)
=== RUN   TestProbeChatRequestStopParse
--- PASS: TestProbeChatRequestStopParse (0.04s)
=== RUN   TestStreamingResponseSkipsDoneAfterErrorChunk
--- PASS: TestStreamingResponseSkipsDoneAfterErrorChunk (0.22s)
=== RUN   TestStreamingResponseSendsDoneOnCleanStream
--- PASS: TestStreamingResponseSendsDoneOnCleanStream (0.20s)
=== RUN   TestSendJSON_StandardBytes
--- PASS: TestSendJSON_StandardBytes (0.00s)
=== RUN   TestSendJSON_Deterministic
--- PASS: TestSendJSON_Deterministic (0.00s)
=== RUN   TestSendJSONWithStatus_StandardBytes
--- PASS: TestSendJSONWithStatus_StandardBytes (0.00s)
=== RUN   TestSendJSON_MarshalError
--- PASS: TestSendJSON_MarshalError (0.00s)
=== RUN   TestSendJSONWithStatus_MarshalError
--- PASS: TestSendJSONWithStatus_MarshalError (0.00s)
=== RUN   TestSendBifrostError_BodylessStatusBecomes502
--- PASS: TestSendBifrostError_BodylessStatusBecomes502 (0.00s)
=== RUN   TestSendBifrostError_NormalStatusPreserved
--- PASS: TestSendBifrostError_NormalStatusPreserved (0.00s)
=== RUN   TestWebhookRoutesRegister
--- PASS: TestWebhookRoutesRegister (0.00s)
=== RUN   TestWebhookHandlerRouteRegistration
    webhooks_test.go:100: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:100
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerRouteRegistration
--- FAIL: TestWebhookHandlerRouteRegistration (0.00s)
=== RUN   TestWebhookHandlerCreateAndGet
    webhooks_test.go:107: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:107
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerCreateAndGet
--- FAIL: TestWebhookHandlerCreateAndGet (0.00s)
=== RUN   TestWebhookHandlerListFilters
    webhooks_test.go:139: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:139
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerListFilters
--- FAIL: TestWebhookHandlerListFilters (0.00s)
=== RUN   TestWebhookHandlerCreateValidation
    webhooks_test.go:199: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:199
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerCreateValidation
--- FAIL: TestWebhookHandlerCreateValidation (0.00s)
=== RUN   TestWebhookHandlerDuplicateName
    webhooks_test.go:216: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:216
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerDuplicateName
--- FAIL: TestWebhookHandlerDuplicateName (0.00s)
=== RUN   TestWebhookHandlerUpdate
    webhooks_test.go:225: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:225
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerUpdate
--- FAIL: TestWebhookHandlerUpdate (0.00s)
=== RUN   TestWebhookHandlerDelete
    webhooks_test.go:246: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:246
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerDelete
--- FAIL: TestWebhookHandlerDelete (0.00s)
=== RUN   TestWebhookHandlerRotateSecret
    webhooks_test.go:262: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:262
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerRotateSecret
--- FAIL: TestWebhookHandlerRotateSecret (0.00s)
=== RUN   TestWebhookHandlerListDeliveries
    webhooks_test.go:293: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:293
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerListDeliveries
--- FAIL: TestWebhookHandlerListDeliveries (0.00s)
=== RUN   TestWebhookHandlerRedeliver
    webhooks_test.go:317: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:317
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerRedeliver
--- FAIL: TestWebhookHandlerRedeliver (0.00s)
=== RUN   TestWebhookHandlerTestDelivery
    webhooks_test.go:349: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:349
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerTestDelivery
--- FAIL: TestWebhookHandlerTestDelivery (0.00s)
=== RUN   TestWebhookHandlerStoreUnavailable
--- PASS: TestWebhookHandlerStoreUnavailable (0.00s)
=== RUN   TestWebhookHandlerPerEndpointTuning
    webhooks_test.go:406: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:406
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerPerEndpointTuning
--- FAIL: TestWebhookHandlerPerEndpointTuning (0.00s)
=== RUN   TestWebhookHandlerHeaders
    webhooks_test.go:430: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-ws-cookie67/transports/bifrost-http/handlers/webhooks_test.go:430
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerHeaders
--- FAIL: TestWebhookHandlerHeaders (0.00s)
=== RUN   TestExtractRealtimeTokenFromAuth
=== PAUSE TestExtractRealtimeTokenFromAuth
=== RUN   TestResolveRealtimeSDPTarget_BaseRouteRequiresProviderPrefix
--- PASS: TestResolveRealtimeSDPTarget_BaseRouteRequiresProviderPrefix (0.00s)
=== RUN   TestResolveRealtimeSDPTarget_BaseRouteNormalizesModel
--- PASS: TestResolveRealtimeSDPTarget_BaseRouteNormalizesModel (0.00s)
=== RUN   TestResolveRealtimeSDPTarget_OpenAIRouteDefaultsProvider
--- PASS: TestResolveRealtimeSDPTarget_OpenAIRouteDefaultsProvider (0.00s)
=== RUN   TestParseCallsWebRTCRequest_RawSDPKeepsGARoute
--- PASS: TestParseCallsWebRTCRequest_RawSDPKeepsGARoute (0.00s)
=== RUN   TestResolveRealtimeSDPTargetDedicatedTranscription
=== PAUSE TestResolveRealtimeSDPTargetDedicatedTranscription
=== RUN   TestResolveRealtimeSDPTargetRootModelPreservesNormalRealtime
=== PAUSE TestResolveRealtimeSDPTargetRootModelPreservesNormalRealtime
=== RUN   TestPinRealtimeSDPTranscriptionModelPreservesSession
=== PAUSE TestPinRealtimeSDPTranscriptionModelPreservesSession
=== RUN   TestNewRealtimeRelayContextCarriesRequestGrant
--- PASS: TestNewRealtimeRelayContextCarriesRequestGrant (0.00s)
=== RUN   TestNewRealtimeRelayContextCopiesValuesWithoutRequestCancellation
--- PASS: TestNewRealtimeRelayContextCopiesValuesWithoutRequestCancellation (0.00s)
=== RUN   TestParseRealtimeEventPreservesExtraParams
--- PASS: TestParseRealtimeEventPreservesExtraParams (0.01s)
=== RUN   TestExtractRealtimeBearerToken
--- PASS: TestExtractRealtimeBearerToken (0.00s)
=== RUN   TestLookupRealtimeEphemeralKeyMappingKeepsEntryUntilTTLExpiry
=== PAUSE TestLookupRealtimeEphemeralKeyMappingKeepsEntryUntilTTLExpiry
=== RUN   TestCacheRealtimeEphemeralKeyMappingPreservesVirtualKeyAcrossKVReplication
=== PAUSE TestCacheRealtimeEphemeralKeyMappingPreservesVirtualKeyAcrossKVReplication
=== RUN   TestLookupRealtimeEphemeralKeyMapping_BackwardsCompatibleStringValue
=== PAUSE TestLookupRealtimeEphemeralKeyMapping_BackwardsCompatibleStringValue
=== RUN   TestWebRTCRealtimeRelayCloseFinalizesActiveTurnHooks
=== PAUSE TestWebRTCRealtimeRelayCloseFinalizesActiveTurnHooks
=== RUN   TestResolveRealtimeWebRTCKeys_UnmappedEphemeralTokenStaysAnonymous
=== PAUSE TestResolveRealtimeWebRTCKeys_UnmappedEphemeralTokenStaysAnonymous
=== RUN   TestApplyRealtimeEphemeralKeyMapping_RestoresVirtualKeyAndKeyID
=== PAUSE TestApplyRealtimeEphemeralKeyMapping_RestoresVirtualKeyAndKeyID
=== RUN   TestBroadcastNotificationTargetsMatchingRole
--- PASS: TestBroadcastNotificationTargetsMatchingRole (0.21s)
=== RUN   TestWebSocketWriteAfterHandlerReturn
--- PASS: TestWebSocketWriteAfterHandlerReturn (0.26s)
=== RUN   TestWebSocketBroadcastDuringDisconnectStorm
--- PASS: TestWebSocketBroadcastDuringDisconnectStorm (0.10s)
=== RUN   TestWebSocketStopClosesClients
--- PASS: TestWebSocketStopClosesClients (0.00s)
=== RUN   TestMatchesWildcardPattern
=== RUN   TestMatchesWildcardPattern/subdomain_match
--- PASS: TestMatchesWildcardPattern/subdomain_match (0.00s)
=== RUN   TestMatchesWildcardPattern/subdomain_no_match
--- PASS: TestMatchesWildcardPattern/subdomain_no_match (0.00s)
=== RUN   TestMatchesWildcardPattern/scheme-less_no_match_with_scheme_prefix
--- PASS: TestMatchesWildcardPattern/scheme-less_no_match_with_scheme_prefix (0.00s)
=== RUN   TestMatchesWildcardPattern/scheme-less_exact_subdomain
--- PASS: TestMatchesWildcardPattern/scheme-less_exact_subdomain (0.00s)
=== RUN   TestMatchesWildcardPattern/no_nested_subdomain
--- PASS: TestMatchesWildcardPattern/no_nested_subdomain (0.00s)
=== RUN   TestMatchesWildcardPattern/empty_subdomain_no_match
--- PASS: TestMatchesWildcardPattern/empty_subdomain_no_match (0.00s)
=== RUN   TestMatchesWildcardPattern/no_wildcard_in_pattern
--- PASS: TestMatchesWildcardPattern/no_wildcard_in_pattern (0.00s)
=== RUN   TestMatchesWildcardPattern/http_vs_https
--- PASS: TestMatchesWildcardPattern/http_vs_https (0.00s)
=== RUN   TestMatchesWildcardPattern/http_wildcard
--- PASS: TestMatchesWildcardPattern/http_wildcard (0.00s)
=== RUN   TestMatchesWildcardPattern/multiple_wildcards
--- PASS: TestMatchesWildcardPattern/multiple_wildcards (0.00s)
=== RUN   TestMatchesWildcardPattern/empty_origin
--- PASS: TestMatchesWildcardPattern/empty_origin (0.00s)
=== RUN   TestMatchesWildcardPattern/empty_pattern
--- PASS: TestMatchesWildcardPattern/empty_pattern (0.00s)
=== RUN   TestMatchesWildcardPattern/both_empty
--- PASS: TestMatchesWildcardPattern/both_empty (0.00s)
=== RUN   TestMatchesWildcardPattern/no_slash_in_wildcard
--- PASS: TestMatchesWildcardPattern/no_slash_in_wildcard (0.00s)
--- PASS: TestMatchesWildcardPattern (0.00s)
=== RUN   TestMatchesWildcardPattern_CacheConsistency
--- PASS: TestMatchesWildcardPattern_CacheConsistency (0.00s)
=== RUN   TestIsOriginAllowed
=== RUN   TestIsOriginAllowed/localhost_http
--- PASS: TestIsOriginAllowed/localhost_http (0.00s)
=== RUN   TestIsOriginAllowed/localhost_https
--- PASS: TestIsOriginAllowed/localhost_https (0.00s)
=== RUN   TestIsOriginAllowed/127.0.0.1
--- PASS: TestIsOriginAllowed/127.0.0.1 (0.00s)
=== RUN   TestIsOriginAllowed/exact_match
--- PASS: TestIsOriginAllowed/exact_match (0.00s)
=== RUN   TestIsOriginAllowed/exact_no_match
--- PASS: TestIsOriginAllowed/exact_no_match (0.00s)
=== RUN   TestIsOriginAllowed/star_allows_all
--- PASS: TestIsOriginAllowed/star_allows_all (0.00s)
=== RUN   TestIsOriginAllowed/wildcard_match
--- PASS: TestIsOriginAllowed/wildcard_match (0.00s)
=== RUN   TestIsOriginAllowed/wildcard_no_match
--- PASS: TestIsOriginAllowed/wildcard_no_match (0.00s)
=== RUN   TestIsOriginAllowed/exact_before_wildcard
--- PASS: TestIsOriginAllowed/exact_before_wildcard (0.00s)
=== RUN   TestIsOriginAllowed/wildcard_after_exact_miss
--- PASS: TestIsOriginAllowed/wildcard_after_exact_miss (0.00s)
=== RUN   TestIsOriginAllowed/empty_origins
--- PASS: TestIsOriginAllowed/empty_origins (0.00s)
=== RUN   TestIsOriginAllowed/empty_origins_slice
--- PASS: TestIsOriginAllowed/empty_origins_slice (0.00s)
--- PASS: TestIsOriginAllowed (0.00s)
=== RUN   TestDiscoverRealtimeTranscriptionModel
=== PAUSE TestDiscoverRealtimeTranscriptionModel
=== RUN   TestPinRealtimeTranscriptionModel
=== PAUSE TestPinRealtimeTranscriptionModel
=== RUN   TestSnapshotRealtimeMiddlewareValuesWithContext
=== PAUSE TestSnapshotRealtimeMiddlewareValuesWithContext
=== RUN   TestBufferRealtimeTranscriptionBootstrapPreservesFrames
--- PASS: TestBufferRealtimeTranscriptionBootstrapPreservesFrames (0.00s)
=== RUN   TestRealtimeStopHeartbeatWaitsForPingGoroutine
--- PASS: TestRealtimeStopHeartbeatWaitsForPingGoroutine (0.05s)
=== RUN   TestRealtimeStopHeartbeatWithoutStart
--- PASS: TestRealtimeStopHeartbeatWithoutStart (0.00s)
=== RUN   TestRealtimeStopHeartbeatIsIdempotent
--- PASS: TestRealtimeStopHeartbeatIsIdempotent (0.02s)
=== RUN   TestNativeWSUpstreamProducesLLMCallSpan
--- PASS: TestNativeWSUpstreamProducesLLMCallSpan (0.11s)
=== RUN   TestNativeWSUpstreamFailedResponseProducesErroredLLMCallSpan
--- PASS: TestNativeWSUpstreamFailedResponseProducesErroredLLMCallSpan (0.01s)
=== RUN   TestNativeWSUpstreamTimeoutProducesErroredLLMCallSpan
--- PASS: TestNativeWSUpstreamTimeoutProducesErroredLLMCallSpan (1.01s)
=== RUN   TestChatGPTWSHandshakeCookies
--- PASS: TestChatGPTWSHandshakeCookies (0.11s)
=== RUN   TestChatGPTWSHandshakeFailureCleanup
=== RUN   TestChatGPTWSHandshakeFailureCleanup/false
--- PASS: TestChatGPTWSHandshakeFailureCleanup/false (0.00s)
=== RUN   TestChatGPTWSHandshakeFailureCleanup/true
--- PASS: TestChatGPTWSHandshakeFailureCleanup/true (0.00s)
--- PASS: TestChatGPTWSHandshakeFailureCleanup (0.00s)
=== RUN   TestChatGPTWSHandshakeAdmission
--- PASS: TestChatGPTWSHandshakeAdmission (0.00s)
=== RUN   TestChatGPTWSMultiTurnRawAndQuota
=== RUN   TestChatGPTWSMultiTurnRawAndQuota/false
--- PASS: TestChatGPTWSMultiTurnRawAndQuota/false (0.00s)
=== RUN   TestChatGPTWSMultiTurnRawAndQuota/true
--- PASS: TestChatGPTWSMultiTurnRawAndQuota/true (0.00s)
--- PASS: TestChatGPTWSMultiTurnRawAndQuota (0.01s)
=== RUN   TestChatGPTWSTargetRejectsPlaintextRemote
--- PASS: TestChatGPTWSTargetRejectsPlaintextRemote (0.00s)
=== RUN   TestChatGPTWSRouteRequiresUpgradeAndCredentials
--- PASS: TestChatGPTWSRouteRequiresUpgradeAndCredentials (0.00s)
=== RUN   TestChatGPTWSDisconnectDuringTurnFinalizesHooks
--- PASS: TestChatGPTWSDisconnectDuringTurnFinalizesHooks (0.00s)
=== RUN   TestResolveWSStreamIdleTimeoutUsesProviderOverride
--- PASS: TestResolveWSStreamIdleTimeoutUsesProviderOverride (0.00s)
=== RUN   TestResolveWSStreamIdleTimeoutFallsBackToDefault
--- PASS: TestResolveWSStreamIdleTimeoutFallsBackToDefault (0.00s)
=== RUN   TestIsWSReadTimeout
--- PASS: TestIsWSReadTimeout (0.00s)
=== RUN   TestNewBifrostError
--- PASS: TestNewBifrostError (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_BaggageSessionIDSetsGrouping
--- PASS: TestCreateBifrostContextFromAuth_BaggageSessionIDSetsGrouping (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_SettlesIdentity
=== RUN   TestCreateBifrostContextFromAuth_SettlesIdentity/virtual_key_header
--- PASS: TestCreateBifrostContextFromAuth_SettlesIdentity/virtual_key_header (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_SettlesIdentity/bearer_virtual_key
--- PASS: TestCreateBifrostContextFromAuth_SettlesIdentity/bearer_virtual_key (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_SettlesIdentity/provider_key_bearer
--- PASS: TestCreateBifrostContextFromAuth_SettlesIdentity/provider_key_bearer (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_SettlesIdentity/no_credential
--- PASS: TestCreateBifrostContextFromAuth_SettlesIdentity/no_credential (0.00s)
--- PASS: TestCreateBifrostContextFromAuth_SettlesIdentity (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_EmptyBaggageSessionIDIgnored
--- PASS: TestCreateBifrostContextFromAuth_EmptyBaggageSessionIDIgnored (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_ForwardsPrefixedHeaders
--- PASS: TestCreateBifrostContextFromAuth_ForwardsPrefixedHeaders (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_AppliesHeaderFilterAndDirectAllowlist
--- PASS: TestCreateBifrostContextFromAuth_AppliesHeaderFilterAndDirectAllowlist (0.00s)
=== RUN   TestCaptureAuthHeaders_PreservesDuplicateHeaderValues
--- PASS: TestCaptureAuthHeaders_PreservesDuplicateHeaderValues (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_PreservesMultipleForwardedHeaderValues
--- PASS: TestCreateBifrostContextFromAuth_PreservesMultipleForwardedHeaderValues (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_BlocksWebSocketHandshakeForwardedHeaders
--- PASS: TestCreateBifrostContextFromAuth_BlocksWebSocketHandshakeForwardedHeaders (0.00s)
=== RUN   TestMergeWebSocketHeaders_ForwardedHeadersOverrideProviderHeadersAndPreserveValues
--- PASS: TestMergeWebSocketHeaders_ForwardedHeadersOverrideProviderHeadersAndPreserveValues (0.00s)
=== RUN   TestHasWebSocketForwardedHeaders
--- PASS: TestHasWebSocketForwardedHeaders (0.00s)
=== RUN   TestSignedWSTicketValidatesAcrossStores
--- PASS: TestSignedWSTicketValidatesAcrossStores (0.00s)
=== RUN   TestSignedWSTicketRejectsWrongKey
--- PASS: TestSignedWSTicketRejectsWrongKey (0.00s)
=== RUN   TestSignedWSTicketRejectsExpiredTicket
--- PASS: TestSignedWSTicketRejectsExpiredTicket (0.00s)
=== RUN   TestSignedWSTicketRejectsMalformedTicket
--- PASS: TestSignedWSTicketRejectsMalformedTicket (0.00s)
=== RUN   TestSignedWSTicketRejectsTamperedTicket
--- PASS: TestSignedWSTicketRejectsTamperedTicket (0.00s)
=== RUN   TestLegacyWSTicketRemainsSingleUse
--- PASS: TestLegacyWSTicketRemainsSingleUse (0.00s)
=== RUN   TestSignedWSTicketDoesNotExposeSessionToken
--- PASS: TestSignedWSTicketDoesNotExposeSessionToken (0.00s)
=== RUN   TestSignedWSTicketPayloadEncryptionRejectsShortPayload
--- PASS: TestSignedWSTicketPayloadEncryptionRejectsShortPayload (0.00s)
=== RUN   TestSignedWSTicketNonceIsHex
--- PASS: TestSignedWSTicketNonceIsHex (0.00s)
=== CONT  TestResolveRealtimeClientSecretTarget
=== CONT  TestPinRealtimeTranscriptionModel
=== RUN   TestResolveRealtimeClientSecretTarget/base_route_with_session_model
=== PAUSE TestResolveRealtimeClientSecretTarget/base_route_with_session_model
=== CONT  TestBuildRealtimeTurnPreRequestPreservesEmptyResponseTools
=== CONT  TestBuildRealtimeTurnPostResponseModelsTranscriptionAsOutput
=== CONT  TestDiscoverRealtimeTranscriptionModel
=== CONT  TestResolveRealtimeWebRTCKeys_UnmappedEphemeralTokenStaysAnonymous
--- PASS: TestBuildRealtimeTurnPostResponseModelsTranscriptionAsOutput (0.00s)
=== CONT  TestResolveRealtimeSDPTargetRootModelPreservesNormalRealtime
--- PASS: TestBuildRealtimeTurnPreRequestPreservesEmptyResponseTools (0.00s)
=== CONT  TestResolveRealtimeSDPTargetDedicatedTranscription
=== RUN   TestPinRealtimeTranscriptionModel/transcription_session_uses_pinned_alias-resolved_model
=== CONT  TestWebRTCRealtimeRelayCloseFinalizesActiveTurnHooks
=== CONT  TestLookupRealtimeEphemeralKeyMapping_BackwardsCompatibleStringValue
=== CONT  TestCacheRealtimeEphemeralKeyMappingPreservesVirtualKeyAcrossKVReplication
--- PASS: TestResolveRealtimeSDPTargetRootModelPreservesNormalRealtime (0.00s)
=== CONT  TestLookupRealtimeEphemeralKeyMappingKeepsEntryUntilTTLExpiry
=== CONT  TestRealtimeMappingVirtualKeyUsesSettledBearerCredential
=== CONT  TestSnapshotRealtimeMiddlewareValuesWithContext
=== CONT  TestPinRealtimeSDPTranscriptionModelPreservesSession
=== CONT  TestRealtimeTurnCompletionContentModelsTranscriptAsOutputOnlyForTranscriptionSessions
=== CONT  TestNewBifrostErrorFromRealtimeErrorCarriesRealtimeMetadata
=== CONT  TestPendingRealtimeToolOutputUpdate
=== CONT  TestFinalizeRealtimeTurnHooksWithErrorCompletesActiveHooks
=== CONT  TestRealtimeMappingVirtualKeyUsesRequestHeader
=== CONT  TestRealtimeTurnFinalEventUsesTranscriptionCompletionOnlyForTranscriptionSessions
=== CONT  TestRealtimeSessionDedupeNestedRawEvents
=== CONT  TestParseRealtimeEphemeralKeyMapping_NestedFallback
=== RUN   TestDiscoverRealtimeTranscriptionModel/nested_transcription_model
=== CONT  TestParseRealtimeEphemeralKeyMapping
=== CONT  TestRewriteGASessionTranscriptionModelNoOpForFullRealtimeSession
=== CONT  TestApplyRealtimeEphemeralKeyMapping_RestoresVirtualKeyAndKeyID
=== CONT  TestBuildRealtimeTurnPreRequestCombinesConversationAndResponseInput
=== RUN   TestResolveRealtimeClientSecretTarget/openai_alias_uses_bare_model
=== CONT  TestGATranscriptionSessionEndToEndThroughFullNormalizationPath
=== PAUSE TestResolveRealtimeClientSecretTarget/openai_alias_uses_bare_model
=== RUN   TestResolveRealtimeClientSecretTarget/base_route_rejects_bare_model
=== PAUSE TestResolveRealtimeClientSecretTarget/base_route_rejects_bare_model
=== RUN   TestResolveRealtimeClientSecretTarget/missing_model
=== PAUSE TestResolveRealtimeClientSecretTarget/missing_model
=== PAUSE TestPinRealtimeTranscriptionModel/transcription_session_uses_pinned_alias-resolved_model
=== RUN   TestPinRealtimeTranscriptionModel/normal_realtime_keeps_input_transcription_model
=== PAUSE TestPinRealtimeTranscriptionModel/normal_realtime_keeps_input_transcription_model
--- PASS: TestSnapshotRealtimeMiddlewareValuesWithContext (0.00s)
--- PASS: TestWebRTCRealtimeRelayCloseFinalizesActiveTurnHooks (0.00s)
--- PASS: TestResolveRealtimeSDPTargetDedicatedTranscription (0.00s)
--- PASS: TestRealtimeTurnCompletionContentModelsTranscriptAsOutputOnlyForTranscriptionSessions (0.00s)
--- PASS: TestNewBifrostErrorFromRealtimeErrorCarriesRealtimeMetadata (0.00s)
--- PASS: TestPendingRealtimeToolOutputUpdate (0.00s)
--- PASS: TestRealtimeMappingVirtualKeyUsesSettledBearerCredential (0.00s)
--- PASS: TestRealtimeMappingVirtualKeyUsesRequestHeader (0.00s)
=== CONT  TestIsJSONContentType
=== CONT  TestRealtimeClientSecretsEvaluateMintingGovernance_PassesContext
=== CONT  TestRealtimeClientSecretsEvaluateMintingGovernance_RequiresAccess
--- PASS: TestRealtimeTurnFinalEventUsesTranscriptionCompletionOnlyForTranscriptionSessions (0.00s)
--- PASS: TestFinalizeRealtimeTurnHooksWithErrorCompletesActiveHooks (0.00s)
--- PASS: TestPinRealtimeSDPTranscriptionModelPreservesSession (0.00s)
--- PASS: TestRealtimeSessionDedupeNestedRawEvents (0.00s)
--- PASS: TestParseRealtimeEphemeralKeyMapping_NestedFallback (0.00s)
=== CONT  TestReplaceAndCacheRealtimeEphemeralToken
=== CONT  TestReplaceAndCacheRealtimeEphemeralTokenRejectsExpiredSecret
=== CONT  TestRewriteGASessionTranscriptionModelAppliesAliasResolution
=== CONT  TestRealtimeMappingVirtualKeyFallsBackToContextValue
=== CONT  TestRealtimeClientSecretsEvaluateMintingGovernance_ContinuesWithoutGovernance
=== CONT  TestParseRealtimeEphemeralKeyMapping_RejectsCompositeShapes
=== CONT  TestBuildRealtimeTurnPreRequestIncludesResponseCreateContent
=== CONT  TestPinRealtimeTranscriptionModel/normal_realtime_keeps_input_transcription_model
=== CONT  TestFullRealtimeSessionWithTranscriptionSiblingEndToEnd
=== PAUSE TestDiscoverRealtimeTranscriptionModel/nested_transcription_model
=== CONT  TestResolveRealtimeClientSecretTarget_NormalizesModel
=== RUN   TestDiscoverRealtimeTranscriptionModel/wrong_event_type
=== RUN   TestResolveRealtimeClientSecretTarget_NormalizesModel/session.model_provider_prefix_stripped
=== PAUSE TestDiscoverRealtimeTranscriptionModel/wrong_event_type
=== RUN   TestResolveRealtimeClientSecretTarget/GA_transcription_session_resolves_model_from_nested_audio_path
=== CONT  TestExtractRealtimeTokenFromAuth
=== PAUSE TestResolveRealtimeClientSecretTarget/GA_transcription_session_resolves_model_from_nested_audio_path
=== RUN   TestExtractRealtimeTokenFromAuth/virtual_key_header
--- PASS: TestRewriteGASessionTranscriptionModelNoOpForFullRealtimeSession (0.00s)
=== PAUSE TestExtractRealtimeTokenFromAuth/virtual_key_header
=== CONT  TestPinRealtimeTranscriptionModel/transcription_session_uses_pinned_alias-resolved_model
=== RUN   TestExtractRealtimeTokenFromAuth/api_key_header
=== PAUSE TestExtractRealtimeTokenFromAuth/api_key_header
=== PAUSE TestResolveRealtimeClientSecretTarget_NormalizesModel/session.model_provider_prefix_stripped
=== RUN   TestResolveRealtimeClientSecretTarget_NormalizesModel/bare_model_unchanged_on_alias_route
=== PAUSE TestResolveRealtimeClientSecretTarget_NormalizesModel/bare_model_unchanged_on_alias_route
=== CONT  TestResolveRealtimeClientSecretTarget_NormalizesModel/session.model_provider_prefix_stripped
=== CONT  TestResolveRealtimeClientSecretTarget/base_route_with_session_model
=== CONT  TestResolveRealtimeClientSecretTarget/GA_transcription_session_resolves_model_from_nested_audio_path
=== CONT  TestResolveRealtimeClientSecretTarget/base_route_rejects_bare_model
--- PASS: TestResolveRealtimeClientSecretTarget/base_route_with_session_model (0.00s)
=== CONT  TestResolveRealtimeClientSecretTarget/missing_model
--- PASS: TestResolveRealtimeClientSecretTarget/missing_model (0.00s)
=== CONT  TestPendingRealtimeInputUpdate
--- PASS: TestResolveRealtimeClientSecretTarget/base_route_rejects_bare_model (0.00s)
--- PASS: TestResolveRealtimeClientSecretTarget/GA_transcription_session_resolves_model_from_nested_audio_path (0.00s)
=== RUN   TestExtractRealtimeTokenFromAuth/google_api_key
=== PAUSE TestExtractRealtimeTokenFromAuth/google_api_key
=== RUN   TestDiscoverRealtimeTranscriptionModel/top-level_model
=== RUN   TestExtractRealtimeTokenFromAuth/authorization_wins
=== PAUSE TestDiscoverRealtimeTranscriptionModel/top-level_model
=== PAUSE TestExtractRealtimeTokenFromAuth/authorization_wins
=== RUN   TestDiscoverRealtimeTranscriptionModel/invalid_JSON
=== CONT  TestResolveRealtimeClientSecretTarget_NormalizesModel/bare_model_unchanged_on_alias_route
=== PAUSE TestDiscoverRealtimeTranscriptionModel/invalid_JSON
=== CONT  TestDiscoverRealtimeTranscriptionModel/nested_transcription_model
=== CONT  TestResolveRealtimeClientSecretTarget/openai_alias_uses_bare_model
=== CONT  TestDiscoverRealtimeTranscriptionModel/top-level_model
=== CONT  TestDiscoverRealtimeTranscriptionModel/wrong_event_type
--- PASS: TestParseRealtimeEphemeralKeyMapping (0.00s)
--- PASS: TestResolveRealtimeClientSecretTarget_NormalizesModel/session.model_provider_prefix_stripped (0.00s)
--- PASS: TestResolveRealtimeClientSecretTarget_NormalizesModel/bare_model_unchanged_on_alias_route (0.00s)
=== RUN   TestExtractRealtimeTokenFromAuth/nil
=== CONT  TestDiscoverRealtimeTranscriptionModel/invalid_JSON
=== PAUSE TestExtractRealtimeTokenFromAuth/nil
--- PASS: TestDiscoverRealtimeTranscriptionModel/invalid_JSON (0.00s)
--- PASS: TestResolveRealtimeClientSecretTarget/openai_alias_uses_bare_model (0.00s)
--- PASS: TestApplyRealtimeEphemeralKeyMapping_RestoresVirtualKeyAndKeyID (0.00s)
--- PASS: TestIsJSONContentType (0.00s)
=== RUN   TestExtractRealtimeTokenFromAuth/authorization_bearer
--- PASS: TestGATranscriptionSessionEndToEndThroughFullNormalizationPath (0.00s)
--- PASS: TestLookupRealtimeEphemeralKeyMapping_BackwardsCompatibleStringValue (0.00s)
=== PAUSE TestExtractRealtimeTokenFromAuth/authorization_bearer
--- PASS: TestResolveRealtimeWebRTCKeys_UnmappedEphemeralTokenStaysAnonymous (0.00s)
--- PASS: TestReplaceAndCacheRealtimeEphemeralTokenRejectsExpiredSecret (0.00s)
--- PASS: TestRealtimeMappingVirtualKeyFallsBackToContextValue (0.00s)
--- PASS: TestRealtimeClientSecretsEvaluateMintingGovernance_ContinuesWithoutGovernance (0.00s)
=== CONT  TestExtractRealtimeTokenFromAuth/google_api_key
=== CONT  TestExtractRealtimeTokenFromAuth/api_key_header
--- PASS: TestExtractRealtimeTokenFromAuth/google_api_key (0.00s)
--- PASS: TestLookupRealtimeEphemeralKeyMappingKeepsEntryUntilTTLExpiry (0.00s)
--- PASS: TestExtractRealtimeTokenFromAuth/api_key_header (0.00s)
--- PASS: TestRewriteGASessionTranscriptionModelAppliesAliasResolution (0.00s)
=== CONT  TestExtractRealtimeTokenFromAuth/authorization_bearer
--- PASS: TestExtractRealtimeTokenFromAuth/authorization_bearer (0.00s)
=== CONT  TestExtractRealtimeTokenFromAuth/nil
--- PASS: TestExtractRealtimeTokenFromAuth/nil (0.00s)
=== CONT  TestExtractRealtimeTokenFromAuth/authorization_wins
--- PASS: TestDiscoverRealtimeTranscriptionModel/top-level_model (0.00s)
--- PASS: TestExtractRealtimeTokenFromAuth/authorization_wins (0.00s)
=== CONT  TestExtractRealtimeTokenFromAuth/virtual_key_header
--- PASS: TestDiscoverRealtimeTranscriptionModel/nested_transcription_model (0.00s)
--- PASS: TestExtractRealtimeTokenFromAuth/virtual_key_header (0.00s)
--- PASS: TestDiscoverRealtimeTranscriptionModel/wrong_event_type (0.00s)
--- PASS: TestRealtimeClientSecretsEvaluateMintingGovernance_PassesContext (0.00s)
--- PASS: TestPinRealtimeTranscriptionModel/transcription_session_uses_pinned_alias-resolved_model (0.00s)
--- PASS: TestRealtimeClientSecretsEvaluateMintingGovernance_RequiresAccess (0.00s)
--- PASS: TestPinRealtimeTranscriptionModel/normal_realtime_keeps_input_transcription_model (0.00s)
--- PASS: TestFullRealtimeSessionWithTranscriptionSiblingEndToEnd (0.00s)
--- PASS: TestParseRealtimeEphemeralKeyMapping_RejectsCompositeShapes (0.00s)
--- PASS: TestReplaceAndCacheRealtimeEphemeralToken (0.00s)
--- PASS: TestPendingRealtimeInputUpdate (0.00s)
--- PASS: TestResolveRealtimeClientSecretTarget_NormalizesModel (0.00s)
--- PASS: TestResolveRealtimeClientSecretTarget (0.00s)
--- PASS: TestExtractRealtimeTokenFromAuth (0.00s)
--- PASS: TestDiscoverRealtimeTranscriptionModel (0.00s)
--- PASS: TestPinRealtimeTranscriptionModel (0.00s)
--- PASS: TestCacheRealtimeEphemeralKeyMappingPreservesVirtualKeyAcrossKVReplication (0.00s)
--- PASS: TestBuildRealtimeTurnPreRequestCombinesConversationAndResponseInput (0.05s)
--- PASS: TestBuildRealtimeTurnPreRequestIncludesResponseCreateContent (0.06s)
FAIL
FAIL	github.com/maximhq/bifrost/transports/bifrost-http/handlers	116.774s

```

### 未修改基底

```text
=== RUN   TestPoolGetAndReturn
--- PASS: TestPoolGetAndReturn (0.03s)
=== RUN   TestPoolMaxIdlePerKey
--- PASS: TestPoolMaxIdlePerKey (0.00s)
=== RUN   TestPoolClose
=== RUN   TestAnthropicRawStreamTextCodecRewritesOnlyTextDelta
--- PASS: TestPoolClose (0.00s)
--- PASS: TestAnthropicRawStreamTextCodecRewritesOnlyTextDelta (0.00s)
=== RUN   TestPoolDialIncludesHandshakeDetails
=== RUN   TestAnthropicRawStreamTextCodecIgnoresNonTextEvents
--- PASS: TestAnthropicRawStreamTextCodecIgnoresNonTextEvents (0.00s)
=== RUN   TestAnthropicRawStreamTextCodecRejectsMalformedEligibleEvents
--- PASS: TestAnthropicRawStreamTextCodecRejectsMalformedEligibleEvents (0.00s)
=== RUN   TestRewriteAnthropicRawRequestBodyRedactsOnlyContentFields
--- PASS: TestRewriteAnthropicRawRequestBodyRedactsOnlyContentFields (0.00s)
--- PASS: TestPoolDialIncludesHandshakeDetails (0.00s)
=== RUN   TestRewriteAnthropicRawRequestBodyRejectsUnmappedLiteral
--- PASS: TestRewriteAnthropicRawRequestBodyRejectsUnmappedLiteral (0.00s)
=== RUN   TestRewriteAnthropicRawRequestBodyRejectsMalformedJSON
=== RUN   TestDialUpstreamCloseDoesNotAffectPoolCapacityCounters
--- PASS: TestRewriteAnthropicRawRequestBodyRejectsMalformedJSON (0.00s)
=== RUN   TestRewriteAnthropicRawRequestBodyRejectsDuplicateKeys
--- PASS: TestRewriteAnthropicRawRequestBodyRejectsDuplicateKeys (0.00s)
=== RUN   TestRewriteAnthropicRawRequestBodyTransformsTargetsDuplicateText
--- PASS: TestRewriteAnthropicRawRequestBodyTransformsTargetsDuplicateText (0.00s)
=== RUN   TestRewriteAnthropicRawRequestBodyTransformsRejectsOriginalMismatch
--- PASS: TestRewriteAnthropicRawRequestBodyTransformsRejectsOriginalMismatch (0.00s)
=== RUN   TestRewriteAnthropicRawRequestBodyTransformsPreservesHistory
--- PASS: TestRewriteAnthropicRawRequestBodyTransformsPreservesHistory (0.00s)
=== RUN   TestRewriteAnthropicRawResponseTransformsTargetsDuplicateText
--- PASS: TestRewriteAnthropicRawResponseTransformsTargetsDuplicateText (0.00s)
=== RUN   TestMustConvertInPassthrough
=== RUN   TestMustConvertInPassthrough/added_message
--- PASS: TestMustConvertInPassthrough/added_message (0.00s)
=== RUN   TestMustConvertInPassthrough/added_advisor
--- PASS: TestMustConvertInPassthrough/added_advisor (0.00s)
=== RUN   TestMustConvertInPassthrough/added_nil_item
--- PASS: TestMustConvertInPassthrough/added_nil_item (0.00s)
=== RUN   TestMustConvertInPassthrough/done_advisor
--- PASS: TestMustConvertInPassthrough/done_advisor (0.00s)
=== RUN   TestMustConvertInPassthrough/done_web_search
--- PASS: TestMustConvertInPassthrough/done_web_search (0.00s)
=== RUN   TestMustConvertInPassthrough/done_web_fetch
--- PASS: TestMustConvertInPassthrough/done_web_fetch (0.00s)
=== RUN   TestMustConvertInPassthrough/done_code_interpreter
--- PASS: TestMustConvertInPassthrough/done_code_interpreter (0.00s)
=== RUN   TestMustConvertInPassthrough/done_computer
--- PASS: TestMustConvertInPassthrough/done_computer (0.00s)
=== RUN   TestMustConvertInPassthrough/done_message
--- PASS: TestMustConvertInPassthrough/done_message (0.00s)
=== RUN   TestMustConvertInPassthrough/done_function_call
--- PASS: TestMustConvertInPassthrough/done_function_call (0.00s)
=== RUN   TestMustConvertInPassthrough/done_mcp_call
--- PASS: TestMustConvertInPassthrough/done_mcp_call (0.00s)
=== RUN   TestMustConvertInPassthrough/done_nil_item
--- PASS: TestMustConvertInPassthrough/done_nil_item (0.00s)
=== RUN   TestMustConvertInPassthrough/web_search_in_progress
--- PASS: TestMustConvertInPassthrough/web_search_in_progress (0.00s)
=== RUN   TestMustConvertInPassthrough/web_search_completed
--- PASS: TestMustConvertInPassthrough/web_search_completed (0.00s)
=== RUN   TestMustConvertInPassthrough/web_fetch_completed
--- PASS: TestMustConvertInPassthrough/web_fetch_completed (0.00s)
=== RUN   TestMustConvertInPassthrough/code_interpreter_code_done
--- PASS: TestMustConvertInPassthrough/code_interpreter_code_done (0.00s)
=== RUN   TestMustConvertInPassthrough/code_interpreter_completed
--- PASS: TestMustConvertInPassthrough/code_interpreter_completed (0.00s)
=== RUN   TestMustConvertInPassthrough/text_delta
--- PASS: TestMustConvertInPassthrough/text_delta (0.00s)
=== RUN   TestMustConvertInPassthrough/function_args_delta
--- PASS: TestMustConvertInPassthrough/function_args_delta (0.00s)
=== RUN   TestMustConvertInPassthrough/content_part_added
--- PASS: TestMustConvertInPassthrough/content_part_added (0.00s)
=== RUN   TestMustConvertInPassthrough/created
--- PASS: TestMustConvertInPassthrough/created (0.00s)
=== RUN   TestMustConvertInPassthrough/completed
--- PASS: TestMustConvertInPassthrough/completed (0.00s)
--- PASS: TestMustConvertInPassthrough (0.00s)
=== RUN   TestAnthropicContainerUploadSurvivesNormalization
--- PASS: TestAnthropicContainerUploadSurvivesNormalization (0.00s)
=== RUN   TestCheckAnthropicPassthrough_OutputConfigEscapeHatch
=== RUN   TestCheckAnthropicPassthrough_OutputConfigEscapeHatch/vertex
--- PASS: TestCheckAnthropicPassthrough_OutputConfigEscapeHatch/vertex (0.00s)
=== RUN   TestCheckAnthropicPassthrough_OutputConfigEscapeHatch/bedrock_mantle
--- PASS: TestCheckAnthropicPassthrough_OutputConfigEscapeHatch/bedrock_mantle (0.00s)
=== RUN   TestCheckAnthropicPassthrough_OutputConfigEscapeHatch/azure
--- PASS: TestCheckAnthropicPassthrough_OutputConfigEscapeHatch/azure (0.00s)
=== RUN   TestCheckAnthropicPassthrough_OutputConfigEscapeHatch/anthropic
--- PASS: TestCheckAnthropicPassthrough_OutputConfigEscapeHatch/anthropic (0.00s)
--- PASS: TestCheckAnthropicPassthrough_OutputConfigEscapeHatch (0.00s)
=== RUN   TestCheckAnthropicPassthrough_VertexBetaKeepsNativeResponseCodec
=== RUN   TestCheckAnthropicPassthrough_VertexBetaKeepsNativeResponseCodec/prompt-caching-scope-2026-01-05
--- PASS: TestCheckAnthropicPassthrough_VertexBetaKeepsNativeResponseCodec/prompt-caching-scope-2026-01-05 (0.00s)
=== RUN   TestCheckAnthropicPassthrough_VertexBetaKeepsNativeResponseCodec/fast-mode-2026-02-01
--- PASS: TestCheckAnthropicPassthrough_VertexBetaKeepsNativeResponseCodec/fast-mode-2026-02-01 (0.00s)
--- PASS: TestCheckAnthropicPassthrough_VertexBetaKeepsNativeResponseCodec (0.00s)
--- PASS: TestDialUpstreamCloseDoesNotAffectPoolCapacityCounters (0.00s)
=== RUN   TestCheckAnthropicPassthroughLegacyCompleteOmitsStreamCodec
=== RUN   TestPoolExpiredConnection
--- PASS: TestCheckAnthropicPassthroughLegacyCompleteOmitsStreamCodec (0.00s)
=== RUN   TestCheckAnthropicPassthrough_OAuthHeaderRouting
--- PASS: TestCheckAnthropicPassthrough_OAuthHeaderRouting (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/mantle_claude-named_alias_to_openai_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/mantle_claude-named_alias_to_openai_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/vertex_claude-named_alias_to_openai_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/vertex_claude-named_alias_to_openai_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/azure_claude-named_alias_to_openai_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/azure_claude-named_alias_to_openai_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/mantle_neutral_alias_to_claude_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/mantle_neutral_alias_to_claude_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/mantle_unaliased_claude_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/mantle_unaliased_claude_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/mantle_unaliased_openai_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/mantle_unaliased_openai_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/azure_deployment_alias_naming_claude_only_in_model_name
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/azure_deployment_alias_naming_claude_only_in_model_name (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/anthropic_provider_always_native
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/anthropic_provider_always_native (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/bedrock_is_never_native
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/bedrock_is_never_native (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_claude_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_claude_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_non-claude_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_non-claude_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_mantle_claude_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_mantle_claude_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_vertex_claude_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_vertex_claude_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_azure_claude_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_azure_claude_model (0.00s)
=== RUN   TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_mantle_openai_model
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily/ingress_mantle_openai_model (0.00s)
--- PASS: TestIsClaudeModel_MultiFamilyProvidersResolveFamily (0.00s)
=== RUN   TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel
=== RUN   TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel/vertex_claude
--- PASS: TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel/vertex_claude (0.00s)
=== RUN   TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel/bedrock_mantle_claude
--- PASS: TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel/bedrock_mantle_claude (0.00s)
=== RUN   TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel/azure_claude
--- PASS: TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel/azure_claude (0.00s)
=== RUN   TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel/bedrock_mantle_openai
--- PASS: TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel/bedrock_mantle_openai (0.00s)
--- PASS: TestCheckAnthropicPassthrough_ProviderPrefixedClaudeModel (0.00s)
=== RUN   TestAnthropicRawArgumentDelta
--- PASS: TestAnthropicRawArgumentDelta (0.00s)
=== RUN   TestAnthropicMessagesRawResponseCarriesExtraFields
=== RUN   TestAnthropicMessagesRawResponseCarriesExtraFields/raw_capture_requested
=== RUN   TestApplyHTTPRequestToCtxKeepsCookieOnce
--- PASS: TestApplyHTTPRequestToCtxKeepsCookieOnce (0.00s)
=== RUN   TestClearCache_OK
--- PASS: TestClearCache_OK (0.00s)
=== RUN   TestClearCache_RejectsEmptyID
--- PASS: TestAnthropicMessagesRawResponseCarriesExtraFields/raw_capture_requested (0.01s)
=== RUN   TestAnthropicMessagesRawResponseCarriesExtraFields/claude_code_passthrough_stays_byte-identical
--- PASS: TestAnthropicMessagesRawResponseCarriesExtraFields/claude_code_passthrough_stays_byte-identical (0.00s)
--- PASS: TestClearCache_RejectsEmptyID (0.00s)
=== RUN   TestAnthropicMessagesRawResponseCarriesExtraFields/no_raw_response_uses_the_converted_shape
=== RUN   TestClearCache_MissingUserValue
--- PASS: TestAnthropicMessagesRawResponseCarriesExtraFields/no_raw_response_uses_the_converted_shape (0.00s)
--- PASS: TestClearCache_MissingUserValue (0.00s)
--- PASS: TestAnthropicMessagesRawResponseCarriesExtraFields (0.01s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases
=== RUN   TestClearCache_PluginErrorReturns500
--- PASS: TestClearCache_PluginErrorReturns500 (0.00s)
=== RUN   TestClearCache_PluginNotLoaded
--- PASS: TestClearCache_PluginNotLoaded (0.00s)
=== RUN   TestClearCacheByKey_OK
--- PASS: TestClearCacheByKey_OK (0.00s)
=== RUN   TestClearCacheByKey_PluginErrorReturns500
--- PASS: TestClearCacheByKey_PluginErrorReturns500 (0.00s)
=== RUN   TestClearCacheByKey_PluginNotLoaded
--- PASS: TestClearCacheByKey_PluginNotLoaded (0.00s)
=== RUN   TestCodexWSIdentityPreservesBusinessBytes
--- PASS: TestCodexWSIdentityPreservesBusinessBytes (0.00s)
=== RUN   TestCodexWSIdentityRejectsInvalidContext
--- PASS: TestCodexWSIdentityRejectsInvalidContext (0.00s)
=== RUN   TestValidateHeaderFilterConfig
=== RUN   TestValidateHeaderFilterConfig/nil_config
--- PASS: TestValidateHeaderFilterConfig/nil_config (0.00s)
=== RUN   TestValidateHeaderFilterConfig/empty_lists
--- PASS: TestValidateHeaderFilterConfig/empty_lists (0.00s)
=== RUN   TestValidateHeaderFilterConfig/empty_allowlist_and_denylist_slices
--- PASS: TestValidateHeaderFilterConfig/empty_allowlist_and_denylist_slices (0.00s)
=== RUN   TestValidateHeaderFilterConfig/valid_allowlist_patterns
--- PASS: TestValidateHeaderFilterConfig/valid_allowlist_patterns (0.00s)
=== RUN   TestValidateHeaderFilterConfig/valid_denylist_patterns
--- PASS: TestValidateHeaderFilterConfig/valid_denylist_patterns (0.00s)
=== RUN   TestValidateHeaderFilterConfig/valid_allowlist_and_denylist_together
--- PASS: TestValidateHeaderFilterConfig/valid_allowlist_and_denylist_together (0.00s)
=== RUN   TestValidateHeaderFilterConfig/whitespace-only_entries_in_allowlist_are_dropped
--- PASS: TestValidateHeaderFilterConfig/whitespace-only_entries_in_allowlist_are_dropped (0.00s)
=== RUN   TestValidateHeaderFilterConfig/whitespace-only_entries_in_denylist_are_dropped
--- PASS: TestValidateHeaderFilterConfig/whitespace-only_entries_in_denylist_are_dropped (0.00s)
=== RUN   TestValidateHeaderFilterConfig/all-empty_allowlist_becomes_effectively_empty
--- PASS: TestValidateHeaderFilterConfig/all-empty_allowlist_becomes_effectively_empty (0.00s)
=== RUN   TestValidateHeaderFilterConfig/security_header_in_allowlist_rejected
--- PASS: TestValidateHeaderFilterConfig/security_header_in_allowlist_rejected (0.00s)
=== RUN   TestValidateHeaderFilterConfig/security_header_in_denylist_rejected
--- PASS: TestValidateHeaderFilterConfig/security_header_in_denylist_rejected (0.00s)
=== RUN   TestValidateHeaderFilterConfig/wildcard_matching_security_header_allowed_(runtime_strips_security_headers)
--- PASS: TestValidateHeaderFilterConfig/wildcard_matching_security_header_allowed_(runtime_strips_security_headers) (0.00s)
=== RUN   TestValidateHeaderFilterConfig/wildcard_prefix_matching_security_headers_allowed_(runtime_strips_security_headers)
--- PASS: TestValidateHeaderFilterConfig/wildcard_prefix_matching_security_headers_allowed_(runtime_strips_security_headers) (0.00s)
=== RUN   TestValidateHeaderFilterConfig/bare_wildcard_in_allowlist_allowed_(runtime_strips_security_headers)
--- PASS: TestValidateHeaderFilterConfig/bare_wildcard_in_allowlist_allowed_(runtime_strips_security_headers) (0.00s)
=== RUN   TestValidateHeaderFilterConfig/wildcard_in_middle_of_pattern_rejected
--- PASS: TestValidateHeaderFilterConfig/wildcard_in_middle_of_pattern_rejected (0.00s)
--- PASS: TestValidateHeaderFilterConfig (0.00s)
=== RUN   TestValidateHeaderFilterConfig_EmptyEntriesDropped
--- PASS: TestValidateHeaderFilterConfig_EmptyEntriesDropped (0.00s)
=== RUN   TestValidateHeaderFilterConfig_EmptyConfigStillForwardsHeaders
--- PASS: TestValidateHeaderFilterConfig_EmptyConfigStillForwardsHeaders (0.00s)
=== RUN   TestGetPasswordPolicyFailures
=== RUN   TestGetPasswordPolicyFailures/valid_password
--- PASS: TestGetPasswordPolicyFailures/valid_password (0.00s)
=== RUN   TestGetPasswordPolicyFailures/missing_all_requirements
--- PASS: TestGetPasswordPolicyFailures/missing_all_requirements (0.00s)
=== RUN   TestGetPasswordPolicyFailures/missing_character_classes
--- PASS: TestGetPasswordPolicyFailures/missing_character_classes (0.00s)
--- PASS: TestGetPasswordPolicyFailures (0.00s)
=== RUN   TestValidateGlobalToolSyncIntervalMinutes
--- PASS: TestValidateGlobalToolSyncIntervalMinutes (0.00s)
=== RUN   TestUpdateConfig_PersistsVKRotationCooldown
    configvkrotationcooldown_test.go:47: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/configvkrotationcooldown_test.go:47
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestUpdateConfig_PersistsVKRotationCooldown
--- FAIL: TestUpdateConfig_PersistsVKRotationCooldown (0.00s)
=== RUN   TestVirtualKeyBudgetOverrideLifecycle
    governance_test.go:235: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/governance_test.go:235
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestVirtualKeyBudgetOverrideLifecycle
--- FAIL: TestVirtualKeyBudgetOverrideLifecycle (0.00s)
=== RUN   TestVirtualKeyBudgetOverrideRejectsDirectMirrorBudget
    governance_test.go:313: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/governance_test.go:313
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestVirtualKeyBudgetOverrideRejectsDirectMirrorBudget
--- FAIL: TestVirtualKeyBudgetOverrideRejectsDirectMirrorBudget (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdatePreservesOmittedAssociation
--- PASS: TestApplyVirtualKeyOwnershipUpdatePreservesOmittedAssociation (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation
=== RUN   TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/set_team_clears_customer
--- PASS: TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/set_team_clears_customer (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/set_customer_clears_team
--- PASS: TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/set_customer_clears_team (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/null_team_clears_both
--- PASS: TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/null_team_clears_both (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/empty_customer_clears_both
--- PASS: TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/empty_customer_clears_both (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/null_team_and_customer_clears_both
--- PASS: TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/null_team_and_customer_clears_both (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/empty_team_and_customer_clears_both
--- PASS: TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/empty_team_and_customer_clears_both (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/team_with_null_customer_sets_team
--- PASS: TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation/team_with_null_customer_sets_team (0.00s)
--- PASS: TestApplyVirtualKeyOwnershipUpdateSwitchesAndClearsAssociation (0.00s)
=== RUN   TestApplyVirtualKeyOwnershipUpdateRejectsDualAssociation
--- PASS: TestApplyVirtualKeyOwnershipUpdateRejectsDualAssociation (0.00s)
=== RUN   TestFindExistingBudgetPrefersIDOverResetDuration
--- PASS: TestFindExistingBudgetPrefersIDOverResetDuration (0.00s)
=== RUN   TestFindExistingBudgetRejectsUnknownID
--- PASS: TestFindExistingBudgetRejectsUnknownID (0.00s)
=== RUN   TestBudgetFrequencyReplaceInheritsUsageFromOriginalBudgets
--- PASS: TestBudgetFrequencyReplaceInheritsUsageFromOriginalBudgets (0.00s)
=== RUN   TestBudgetLookupConsumesMatchedRowsForDurationSwap
--- PASS: TestBudgetLookupConsumesMatchedRowsForDurationSwap (0.00s)
=== RUN   TestResetBudgetUsageIfRequested
--- PASS: TestResetBudgetUsageIfRequested (0.00s)
=== RUN   TestTeamBudgetFrequencyChangePreservesUsageWhenRequested
--- PASS: TestTeamBudgetFrequencyChangePreservesUsageWhenRequested (0.00s)
=== RUN   TestVirtualKeyBudgetFrequencyChangePreservesUsageWhenRequested
--- PASS: TestVirtualKeyBudgetFrequencyChangePreservesUsageWhenRequested (0.00s)
=== RUN   TestProviderBudgetFrequencyChangePreservesUsageWhenRequested
--- PASS: TestProviderBudgetFrequencyChangePreservesUsageWhenRequested (0.00s)
=== RUN   TestBudgetFrequencyChangeResetsUsageWhenRequested
--- PASS: TestBudgetFrequencyChangeResetsUsageWhenRequested (0.00s)
=== RUN   TestExistingVirtualKeyBudgetLoweredBelowPreservedUsageIsAllowed
--- PASS: TestExistingVirtualKeyBudgetLoweredBelowPreservedUsageIsAllowed (0.00s)
=== RUN   TestExistingProviderBudgetLoweredBelowPreservedUsageIsAllowed
--- PASS: TestExistingProviderBudgetLoweredBelowPreservedUsageIsAllowed (0.00s)
=== RUN   TestExistingBudgetLoweredBelowUsageSucceedsWhenResetRequested
--- PASS: TestExistingBudgetLoweredBelowUsageSucceedsWhenResetRequested (0.00s)
=== RUN   TestNewVirtualKeyBudgetInheritsClosestShorterUsage
--- PASS: TestNewVirtualKeyBudgetInheritsClosestShorterUsage (0.00s)
=== RUN   TestNewProviderBudgetInheritsClosestShorterUsage
--- PASS: TestNewProviderBudgetInheritsClosestShorterUsage (0.00s)
=== RUN   TestNewVirtualKeyBudgetInheritanceAboveLimitIsAllowed
--- PASS: TestNewVirtualKeyBudgetInheritanceAboveLimitIsAllowed (0.00s)
=== RUN   TestNewProviderBudgetInheritanceAtLimitIsAllowed
--- PASS: TestNewProviderBudgetInheritanceAtLimitIsAllowed (0.00s)
=== RUN   TestNewShorterBudgetDoesNotInheritFromLongerUsage
--- PASS: TestNewShorterBudgetDoesNotInheritFromLongerUsage (0.00s)
=== RUN   TestNewBudgetInheritsClosestShorterUsage
--- PASS: TestNewBudgetInheritsClosestShorterUsage (0.00s)
=== RUN   TestNewLongerBudgetDoesNotInheritUsageWhenResetRequested
--- PASS: TestNewLongerBudgetDoesNotInheritUsageWhenResetRequested (0.00s)
=== RUN   TestRotateVirtualKey_OnlyChangesValueAndReloads
--- PASS: TestRotateVirtualKey_OnlyChangesValueAndReloads (0.00s)
=== RUN   TestRotateVirtualKey_CooldownStoresPreviousValue
--- PASS: TestRotateVirtualKey_CooldownStoresPreviousValue (0.00s)
=== RUN   TestRotateVirtualKey_ZeroCooldownClearsPreviousValue
--- PASS: TestRotateVirtualKey_ZeroCooldownClearsPreviousValue (0.00s)
=== RUN   TestRotateVirtualKey_DefaultUnsetCooldownRevokesImmediately
--- PASS: TestRotateVirtualKey_DefaultUnsetCooldownRevokesImmediately (0.00s)
=== RUN   TestRotateVirtualKey_NoClientConfigRevokesImmediately
--- PASS: TestRotateVirtualKey_NoClientConfigRevokesImmediately (0.00s)
=== RUN   TestRotateVirtualKey_ClientConfigErrorRevokesImmediately
--- PASS: TestRotateVirtualKey_ClientConfigErrorRevokesImmediately (0.00s)
=== RUN   TestRotateVirtualKeys_BulkDefaultCooldownRevokesImmediately
--- PASS: TestRotateVirtualKeys_BulkDefaultCooldownRevokesImmediately (0.00s)
=== RUN   TestRotateVirtualKey_NotFound
--- PASS: TestRotateVirtualKey_NotFound (0.00s)
=== RUN   TestRotateVirtualKey_UpdateFailureDoesNotReload
--- PASS: TestRotateVirtualKey_UpdateFailureDoesNotReload (0.00s)
=== RUN   TestRotateVirtualKey_ReloadFailureReturnsErrorAfterUpdate
--- PASS: TestRotateVirtualKey_ReloadFailureReturnsErrorAfterUpdate (0.00s)
=== RUN   TestRotateVirtualKeys_PartialSuccess
--- PASS: TestRotateVirtualKeys_PartialSuccess (0.00s)
=== RUN   TestRotateVirtualKeys_RejectsInvalidRequests
=== RUN   TestRotateVirtualKeys_RejectsInvalidRequests/invalid_JSON
--- PASS: TestRotateVirtualKeys_RejectsInvalidRequests/invalid_JSON (0.00s)
=== RUN   TestRotateVirtualKeys_RejectsInvalidRequests/empty_IDs
--- PASS: TestRotateVirtualKeys_RejectsInvalidRequests/empty_IDs (0.00s)
=== RUN   TestRotateVirtualKeys_RejectsInvalidRequests/blank_ID
--- PASS: TestRotateVirtualKeys_RejectsInvalidRequests/blank_ID (0.00s)
--- PASS: TestRotateVirtualKeys_RejectsInvalidRequests (0.00s)
=== RUN   TestRotateVirtualKeys_TrimsAndDeduplicatesIDs
--- PASS: TestRotateVirtualKeys_TrimsAndDeduplicatesIDs (0.00s)
=== RUN   TestRotateVirtualKeys_AllFailuresReturnsServerError
--- PASS: TestRotateVirtualKeys_AllFailuresReturnsServerError (0.00s)
=== RUN   TestGetVirtualKeyQuota_HydratesBudgetsFromModelConfigs
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/continue_is_refused
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/continue_is_refused (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/create_falls_through
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/create_falls_through (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/no_thread_falls_through
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/no_thread_falls_through (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/malformed_thread_falls_through
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/malformed_thread_falls_through (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/nil_extra_params_falls_through
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/nil_extra_params_falls_through (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/large_payload_continue_refused_from_metadata
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/large_payload_continue_refused_from_metadata (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/large_payload_create_falls_through
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/large_payload_create_falls_through (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/large_payload_without_thread_metadata_falls_through
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/large_payload_without_thread_metadata_falls_through (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/parsed_thread_wins_over_metadata
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/parsed_thread_wins_over_metadata (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/wrong_request_type_falls_through
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/wrong_request_type_falls_through (0.00s)
=== RUN   TestAnthropicRefuseThreadContinue_EdgeCases/count_tokens_path_never_refused
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases/count_tokens_path_never_refused (0.00s)
--- PASS: TestAnthropicRefuseThreadContinue_EdgeCases (0.02s)
=== RUN   TestGenericRouter_MarkDeprecatedListModelsResponseUsesCatalog
--- PASS: TestGetVirtualKeyQuota_HydratesBudgetsFromModelConfigs (0.02s)
=== RUN   TestGetVirtualKeyQuota_ExternalResolverReplacesWithAccessProfileBudgets
--- PASS: TestGetVirtualKeyQuota_ExternalResolverReplacesWithAccessProfileBudgets (0.00s)
=== RUN   TestGetVirtualKeyQuota_ExternalResolverRateLimitOnly
--- PASS: TestGetVirtualKeyQuota_ExternalResolverRateLimitOnly (0.00s)
=== RUN   TestApplyExternalBudgets_RateLimitOnlyDropsNativeBudgets
--- PASS: TestApplyExternalBudgets_RateLimitOnlyDropsNativeBudgets (0.00s)
=== RUN   TestApplyExternalBudgets_ManagedWithNoGovernanceFlagsAndClears
--- PASS: TestApplyExternalBudgets_ManagedWithNoGovernanceFlagsAndClears (0.00s)
=== RUN   TestGetVirtualKeyQuota_ExternalResolverErrorFailsClosed
--- PASS: TestGetVirtualKeyQuota_ExternalResolverErrorFailsClosed (0.00s)
=== RUN   TestGetVirtualKeyQuota_NoGovernanceReturnsEmpty
--- PASS: TestGetVirtualKeyQuota_NoGovernanceReturnsEmpty (0.00s)
=== RUN   TestGetVirtualKeyQuota_MissingHeaderReturns401
--- PASS: TestGetVirtualKeyQuota_MissingHeaderReturns401 (0.00s)
=== RUN   TestGetVirtualKeyQuota_NotFoundReturns401
--- PASS: TestGetVirtualKeyQuota_NotFoundReturns401 (0.00s)
=== RUN   TestGetVirtualKeyQuota_ModelConfigLoadErrorFailsClosed
--- PASS: TestGetVirtualKeyQuota_ModelConfigLoadErrorFailsClosed (0.00s)
=== RUN   TestGetVirtualKeyQuota_RankingsErrorFailsClosed
--- PASS: TestGetVirtualKeyQuota_RankingsErrorFailsClosed (0.00s)
=== RUN   TestGetVirtualKeyQuota_EndToEndWithRealStore
    governance_test.go:2309: failed to create config store: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
--- FAIL: TestGetVirtualKeyQuota_EndToEndWithRealStore (0.00s)
=== RUN   TestGetVirtualKeyQuota_GraceValueWithRealStore
    governance_test.go:2499: failed to create config store: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
--- FAIL: TestGetVirtualKeyQuota_GraceValueWithRealStore (0.00s)
=== RUN   TestGetVirtualKeyQuota_ExpiredGraceValueUnauthorized
    governance_test.go:2524: failed to create config store: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
--- FAIL: TestGetVirtualKeyQuota_ExpiredGraceValueUnauthorized (0.00s)
=== RUN   TestGetVirtualKeyQuota_ExpiredVirtualKeyRejected
--- PASS: TestGetVirtualKeyQuota_ExpiredVirtualKeyRejected (0.00s)
=== RUN   TestGetVirtualKeyQuota_UnexpiredVirtualKeyAllowed
--- PASS: TestGetVirtualKeyQuota_UnexpiredVirtualKeyAllowed (0.00s)
=== RUN   TestGetVirtualKeyQuota_InactiveVirtualKeyStillReadsQuota
--- PASS: TestGetVirtualKeyQuota_InactiveVirtualKeyStillReadsQuota (0.00s)
=== RUN   TestGetVirtualKeyQuota_WindowClampedToBudgetCreation
    governance_test.go:2646: failed to create config store: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
--- FAIL: TestGetVirtualKeyQuota_WindowClampedToBudgetCreation (0.00s)
=== RUN   TestGetVirtualKeys_PaginatedEndpoint_ResponseShape
--- PASS: TestGetVirtualKeys_PaginatedEndpoint_ResponseShape (0.00s)
=== RUN   TestGetVirtualKeys_PaginatedEndpoint_QueryParams
=== RUN   TestGetVirtualKeys_PaginatedEndpoint_QueryParams/explicit_limit_and_offset
--- PASS: TestGetVirtualKeys_PaginatedEndpoint_QueryParams/explicit_limit_and_offset (0.00s)
=== RUN   TestGetVirtualKeys_PaginatedEndpoint_QueryParams/no_params_uses_defaults
--- PASS: TestGetVirtualKeys_PaginatedEndpoint_QueryParams/no_params_uses_defaults (0.00s)
--- PASS: TestGetVirtualKeys_PaginatedEndpoint_QueryParams (0.00s)
=== RUN   TestGetVirtualKeys_FromMemoryUsesGovernanceData
--- PASS: TestGetVirtualKeys_FromMemoryUsesGovernanceData (0.00s)
=== RUN   TestGetVirtualKeys_FromMemoryTakesPrecedenceOverLimit
--- PASS: TestGetVirtualKeys_FromMemoryTakesPrecedenceOverLimit (0.00s)
=== RUN   TestGetVirtualKeys_FromMemoryRejectsUserFilter
--- PASS: TestGetVirtualKeys_FromMemoryRejectsUserFilter (0.00s)
=== RUN   TestGetVirtualKeys_FromMemoryIgnoresOtherFilters
--- PASS: TestGetVirtualKeys_FromMemoryIgnoresOtherFilters (0.00s)
=== RUN   TestBudgetRemovalRequestDetection
=== RUN   TestBudgetRemovalRequestDetection/nil_request_is_not_removal
--- PASS: TestBudgetRemovalRequestDetection/nil_request_is_not_removal (0.00s)
=== RUN   TestBudgetRemovalRequestDetection/empty_object_is_removal
--- PASS: TestBudgetRemovalRequestDetection/empty_object_is_removal (0.00s)
=== RUN   TestBudgetRemovalRequestDetection/max_limit_present_is_not_removal
--- PASS: TestBudgetRemovalRequestDetection/max_limit_present_is_not_removal (0.00s)
=== RUN   TestBudgetRemovalRequestDetection/reset_duration_only_is_not_removal
--- PASS: TestBudgetRemovalRequestDetection/reset_duration_only_is_not_removal (0.00s)
--- PASS: TestBudgetRemovalRequestDetection (0.00s)
=== RUN   TestRateLimitRemovalRequestDetection
=== RUN   TestRateLimitRemovalRequestDetection/nil_request_is_not_removal
--- PASS: TestRateLimitRemovalRequestDetection/nil_request_is_not_removal (0.00s)
=== RUN   TestRateLimitRemovalRequestDetection/empty_object_is_removal
--- PASS: TestRateLimitRemovalRequestDetection/empty_object_is_removal (0.00s)
=== RUN   TestRateLimitRemovalRequestDetection/token_limit_present_is_not_removal
--- PASS: TestRateLimitRemovalRequestDetection/token_limit_present_is_not_removal (0.00s)
=== RUN   TestRateLimitRemovalRequestDetection/request_limit_present_is_not_removal
--- PASS: TestRateLimitRemovalRequestDetection/request_limit_present_is_not_removal (0.00s)
=== RUN   TestRateLimitRemovalRequestDetection/durations_only_is_not_removal
--- PASS: TestRateLimitRemovalRequestDetection/durations_only_is_not_removal (0.00s)
--- PASS: TestRateLimitRemovalRequestDetection (0.00s)
=== RUN   TestCollectProviderConfigDeleteIDs
=== RUN   TestCollectProviderConfigDeleteIDs/collects_both_IDs
--- PASS: TestCollectProviderConfigDeleteIDs/collects_both_IDs (0.00s)
=== RUN   TestCollectProviderConfigDeleteIDs/appends_to_existing_slices
--- PASS: TestCollectProviderConfigDeleteIDs/appends_to_existing_slices (0.00s)
=== RUN   TestCollectProviderConfigDeleteIDs/ignores_missing_IDs
--- PASS: TestCollectProviderConfigDeleteIDs/ignores_missing_IDs (0.00s)
--- PASS: TestCollectProviderConfigDeleteIDs (0.00s)
=== RUN   TestCoerceLegacyBudget
=== RUN   TestCoerceLegacyBudget/empty_object_→_removal,_returns_empty_slice
--- PASS: TestCoerceLegacyBudget/empty_object_→_removal,_returns_empty_slice (0.00s)
=== RUN   TestCoerceLegacyBudget/both_fields_set,_no_existing_→_new_budget_entry,_no_ID
--- PASS: TestCoerceLegacyBudget/both_fields_set,_no_existing_→_new_budget_entry,_no_ID (0.00s)
=== RUN   TestCoerceLegacyBudget/update_max_limit_only,_existing_budget_→_merges_ID_and_reset_duration
--- PASS: TestCoerceLegacyBudget/update_max_limit_only,_existing_budget_→_merges_ID_and_reset_duration (0.00s)
=== RUN   TestCoerceLegacyBudget/update_reset_duration_only,_existing_budget_→_merges_ID_and_max_limit
--- PASS: TestCoerceLegacyBudget/update_reset_duration_only,_existing_budget_→_merges_ID_and_max_limit (0.00s)
=== RUN   TestCoerceLegacyBudget/max_limit_only,_no_existing_→_cannot_build_valid_budget,_returns_nil
--- PASS: TestCoerceLegacyBudget/max_limit_only,_no_existing_→_cannot_build_valid_budget,_returns_nil (0.00s)
=== RUN   TestCoerceLegacyBudget/reset_duration_only,_no_existing_→_cannot_build_valid_budget,_returns_nil
--- PASS: TestCoerceLegacyBudget/reset_duration_only,_no_existing_→_cannot_build_valid_budget,_returns_nil (0.00s)
--- PASS: TestCoerceLegacyBudget (0.00s)
=== RUN   TestModelConfigToProviderGovernanceNewFields
=== RUN   TestModelConfigToProviderGovernanceNewFields/nil_mc_returns_false
--- PASS: TestModelConfigToProviderGovernanceNewFields/nil_mc_returns_false (0.00s)
=== RUN   TestModelConfigToProviderGovernanceNewFields/wrong_scope_returns_false
--- PASS: TestModelConfigToProviderGovernanceNewFields/wrong_scope_returns_false (0.00s)
=== RUN   TestModelConfigToProviderGovernanceNewFields/no_budgets:_Budget_nil,_Budgets_empty,_CalendarAligned_false
--- PASS: TestModelConfigToProviderGovernanceNewFields/no_budgets:_Budget_nil,_Budgets_empty,_CalendarAligned_false (0.00s)
=== RUN   TestModelConfigToProviderGovernanceNewFields/single_budget:_Budget_points_to_first,_Budgets_has_one_entry
--- PASS: TestModelConfigToProviderGovernanceNewFields/single_budget:_Budget_points_to_first,_Budgets_has_one_entry (0.00s)
=== RUN   TestModelConfigToProviderGovernanceNewFields/multiple_budgets:_Budget_is_first,_Budgets_contains_all
--- PASS: TestModelConfigToProviderGovernanceNewFields/multiple_budgets:_Budget_is_first,_Budgets_contains_all (0.00s)
=== RUN   TestModelConfigToProviderGovernanceNewFields/calendar_aligned_is_propagated
--- PASS: TestModelConfigToProviderGovernanceNewFields/calendar_aligned_is_propagated (0.00s)
=== RUN   TestModelConfigToProviderGovernanceNewFields/Budgets_slice_is_a_copy,_not_a_reference_to_mc.Budgets
--- PASS: TestModelConfigToProviderGovernanceNewFields/Budgets_slice_is_a_copy,_not_a_reference_to_mc.Budgets (0.00s)
--- PASS: TestModelConfigToProviderGovernanceNewFields (0.00s)
=== RUN   TestUpdateProviderGovernance_BudgetMutualExclusion
--- PASS: TestUpdateProviderGovernance_BudgetMutualExclusion (0.00s)
=== RUN   TestValidateRoutingFallbacks
=== RUN   TestValidateRoutingFallbacks/nil
--- PASS: TestValidateRoutingFallbacks/nil (0.00s)
=== RUN   TestValidateRoutingFallbacks/empty
--- PASS: TestValidateRoutingFallbacks/empty (0.00s)
=== RUN   TestValidateRoutingFallbacks/provider_model
--- PASS: TestValidateRoutingFallbacks/provider_model (0.00s)
=== RUN   TestValidateRoutingFallbacks/provider_slash_incoming_model
--- PASS: TestValidateRoutingFallbacks/provider_slash_incoming_model (0.00s)
=== RUN   TestValidateRoutingFallbacks/bare_known_provider_name_rejected
--- PASS: TestValidateRoutingFallbacks/bare_known_provider_name_rejected (0.00s)
=== RUN   TestValidateRoutingFallbacks/bare_model_rejected
--- PASS: TestValidateRoutingFallbacks/bare_model_rejected (0.00s)
=== RUN   TestValidateRoutingFallbacks/empty_element
--- PASS: TestValidateRoutingFallbacks/empty_element (0.00s)
=== RUN   TestValidateRoutingFallbacks/huggingface_namespace_not_a_provider_prefix
--- PASS: TestValidateRoutingFallbacks/huggingface_namespace_not_a_provider_prefix (0.00s)
--- PASS: TestValidateRoutingFallbacks (0.00s)
=== RUN   TestCreateCustomer_CalendarAligned_SnapsBudgetLastReset
--- PASS: TestCreateCustomer_CalendarAligned_SnapsBudgetLastReset (0.00s)
=== RUN   TestCreateCustomer_CalendarAligned_False
--- PASS: TestCreateCustomer_CalendarAligned_False (0.00s)
=== RUN   TestUpdateCustomer_CalendarAligned_DoesNotTouchBudgets
--- PASS: TestUpdateCustomer_CalendarAligned_DoesNotTouchBudgets (0.01s)
=== RUN   TestUpdateCustomer_CalendarAligned_NoSnapWhenAlreadyEnabled
--- PASS: TestUpdateCustomer_CalendarAligned_NoSnapWhenAlreadyEnabled (0.00s)
=== RUN   TestApplyVKGovernanceFromModelConfigs_PreservesDirectlyAttachedBudget
--- PASS: TestApplyVKGovernanceFromModelConfigs_PreservesDirectlyAttachedBudget (0.00s)
=== RUN   TestApplyVKGovernanceFromModelConfigs_OverlaysModelConfigGovernance
--- PASS: TestApplyVKGovernanceFromModelConfigs_OverlaysModelConfigGovernance (0.00s)
=== RUN   TestProviderGovernance_DecodesEncodedProviderName
    governance_test.go:3752: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/governance_test.go:3752
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestProviderGovernance_DecodesEncodedProviderName
--- FAIL: TestProviderGovernance_DecodesEncodedProviderName (0.00s)
=== RUN   TestProviderGovernance_UnknownProviderStill404
    governance_test.go:3802: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/governance_test.go:3802
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestProviderGovernance_UnknownProviderStill404
--- FAIL: TestProviderGovernance_UnknownProviderStill404 (0.00s)
=== RUN   TestUpdateProviderGovernance_AdoptsReconciledBudgetsNotStaleOnes
    governance_test.go:3846: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/governance_test.go:3846
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestUpdateProviderGovernance_AdoptsReconciledBudgetsNotStaleOnes
--- FAIL: TestUpdateProviderGovernance_AdoptsReconciledBudgetsNotStaleOnes (0.00s)
=== RUN   TestProviderGovernance_MalformedEncodingReturns400
    governance_test.go:3895: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/governance_test.go:3895
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestProviderGovernance_MalformedEncodingReturns400
--- FAIL: TestProviderGovernance_MalformedEncodingReturns400 (0.00s)
=== RUN   TestTeam_DecodesEncodedTeamID
    governance_test.go:3933: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/governance_test.go:3933
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestTeam_DecodesEncodedTeamID
--- FAIL: TestTeam_DecodesEncodedTeamID (0.00s)
=== RUN   TestTeam_MalformedEncodingReturns400
    governance_test.go:4013: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/governance_test.go:4013
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestTeam_MalformedEncodingReturns400
--- FAIL: TestTeam_MalformedEncodingReturns400 (0.00s)
=== RUN   TestValidateBudgetResetConfig
=== RUN   TestValidateBudgetResetConfig/quarterly_with_no_config_is_valid
--- PASS: TestValidateBudgetResetConfig/quarterly_with_no_config_is_valid (0.00s)
=== RUN   TestValidateBudgetResetConfig/quarterly_with_a_fiscal_start_is_valid
--- PASS: TestValidateBudgetResetConfig/quarterly_with_a_fiscal_start_is_valid (0.00s)
=== RUN   TestValidateBudgetResetConfig/monthly_with_a_quarter_config_is_rejected
--- PASS: TestValidateBudgetResetConfig/monthly_with_a_quarter_config_is_rejected (0.00s)
=== RUN   TestValidateBudgetResetConfig/hourly_with_a_quarter_config_is_rejected
--- PASS: TestValidateBudgetResetConfig/hourly_with_a_quarter_config_is_rejected (0.00s)
=== RUN   TestValidateBudgetResetConfig/month_above_december_is_rejected
--- PASS: TestValidateBudgetResetConfig/month_above_december_is_rejected (0.00s)
=== RUN   TestValidateBudgetResetConfig/negative_month_is_rejected
--- PASS: TestValidateBudgetResetConfig/negative_month_is_rejected (0.00s)
=== RUN   TestValidateBudgetResetConfig/zero_month_means_unset_and_is_accepted
--- PASS: TestValidateBudgetResetConfig/zero_month_means_unset_and_is_accepted (0.00s)
--- PASS: TestValidateBudgetResetConfig (0.00s)
=== RUN   TestNewBudgetFromRequestCarriesResetConfig
--- PASS: TestNewBudgetFromRequestCarriesResetConfig (0.00s)
=== RUN   TestNewBudgetFromRequestSnapsToFiscalQuarter
--- PASS: TestNewBudgetFromRequestSnapsToFiscalQuarter (0.00s)
=== RUN   TestApplyResetConfigToExistingBudgetIsConfigOnly
--- PASS: TestApplyResetConfigToExistingBudgetIsConfigOnly (0.00s)
=== RUN   TestApplyResetConfigToExistingBudgetClearsDefinition
--- PASS: TestApplyResetConfigToExistingBudgetClearsDefinition (0.00s)
=== RUN   TestQuarterStartChangeMakesBudgetDue
--- PASS: TestQuarterStartChangeMakesBudgetDue (0.00s)
=== RUN   TestBudgetLastResetUsesBudgetQuarterStart
--- PASS: TestBudgetLastResetUsesBudgetQuarterStart (0.00s)
=== RUN   TestApplyAssignees
=== RUN   TestApplyAssignees/fills_in_assignees_in_one_batched_call
--- PASS: TestApplyAssignees/fills_in_assignees_in_one_batched_call (0.00s)
=== RUN   TestApplyAssignees/no-ops_without_a_resolver
--- PASS: TestApplyAssignees/no-ops_without_a_resolver (0.00s)
=== RUN   TestApplyAssignees/degrades_to_no_assignee_when_the_resolver_fails
--- PASS: TestApplyAssignees/degrades_to_no_assignee_when_the_resolver_fails (0.00s)
=== RUN   TestApplyAssignees/skips_the_resolver_for_an_empty_page
--- PASS: TestApplyAssignees/skips_the_resolver_for_an_empty_page (0.00s)
--- PASS: TestApplyAssignees (0.00s)
=== RUN   TestGovernanceRouteOverridesReplaceTeamFamily
--- PASS: TestGovernanceRouteOverridesReplaceTeamFamily (0.00s)
=== RUN   TestGovernanceRouteOverridesDefaultToOSS
--- PASS: TestGovernanceRouteOverridesDefaultToOSS (0.00s)
=== RUN   TestResolveBatchProvider
=== RUN   TestResolveBatchProvider/model_field:_provider+model_parsed
--- PASS: TestResolveBatchProvider/model_field:_provider+model_parsed (0.00s)
=== RUN   TestResolveBatchProvider/no_model,_x-model-provider_header
--- PASS: TestResolveBatchProvider/no_model,_x-model-provider_header (0.00s)
=== RUN   TestResolveBatchProvider/no_model,_?provider=_query_param
--- PASS: TestResolveBatchProvider/no_model,_?provider=_query_param (0.00s)
=== RUN   TestResolveBatchProvider/no_model,_no_provider_→_error
--- PASS: TestResolveBatchProvider/no_model,_no_provider_→_error (0.00s)
--- PASS: TestResolveBatchProvider (0.00s)
=== RUN   TestPrepareImageEditRequest_JSON
--- PASS: TestPrepareImageEditRequest_JSON (0.01s)
=== RUN   TestPrepareImageEditRequest_JSONTypedExtraParams
--- PASS: TestPrepareImageEditRequest_JSONTypedExtraParams (0.00s)
=== RUN   TestPrepareImageEditRequest_JSONInputFieldsAreNotExtraParams
--- PASS: TestPrepareImageEditRequest_JSONInputFieldsAreNotExtraParams (0.00s)
=== RUN   TestPrepareImageEditRequest_JSONBareStringImages
--- PASS: TestPrepareImageEditRequest_JSONBareStringImages (0.00s)
=== RUN   TestPrepareImageEditRequest_JSONBase64Image
--- PASS: TestPrepareImageEditRequest_JSONBase64Image (0.00s)
=== RUN   TestPrepareImageEditRequest_MultipartUnchanged
--- PASS: TestPrepareImageEditRequest_MultipartUnchanged (0.00s)
=== RUN   TestPrepareImageEditRequest_MultipartUpload
--- PASS: TestPrepareImageEditRequest_MultipartUpload (0.00s)
=== RUN   TestPrepareImageEditRequest_Validation
--- PASS: TestPrepareImageEditRequest_Validation (0.00s)
=== RUN   TestPrepareImageEditRequest_EmptyImageEntriesAreDropped
--- PASS: TestPrepareImageEditRequest_EmptyImageEntriesAreDropped (0.00s)
=== RUN   TestPrepareImageEditRequest_Stream
--- PASS: TestPrepareImageEditRequest_Stream (0.00s)
=== RUN   TestApplyListModelsProviderFilterDelegatesToTheModelsManager
--- PASS: TestApplyListModelsProviderFilterDelegatesToTheModelsManager (0.00s)
=== RUN   TestApplyListModelsProviderFilterLeavesFanOutAloneWhenNothingResolved
--- PASS: TestApplyListModelsProviderFilterLeavesFanOutAloneWhenNothingResolved (0.00s)
=== RUN   TestApplyListModelsProviderFilterWithoutModelsManager
--- PASS: TestApplyListModelsProviderFilterWithoutModelsManager (0.00s)
=== RUN   TestPrepareVideoEditRequest_JSON
--- PASS: TestPrepareVideoEditRequest_JSON (0.00s)
=== RUN   TestPrepareVideoEditRequest_ProviderFromVideoID
--- PASS: TestPrepareVideoEditRequest_ProviderFromVideoID (0.00s)
=== RUN   TestPrepareVideoEditRequest_ProviderFallbacks
--- PASS: TestPrepareVideoEditRequest_ProviderFallbacks (0.00s)
=== RUN   TestPrepareVideoEditRequest_BareModelStillResolvesProvider
--- PASS: TestPrepareVideoEditRequest_BareModelStillResolvesProvider (0.00s)
=== RUN   TestPrepareVideoEditRequest_BareModelDefersToRouting
--- PASS: TestPrepareVideoEditRequest_BareModelDefersToRouting (0.00s)
=== RUN   TestVideoIDProviderSuffix
--- PASS: TestVideoIDProviderSuffix (0.00s)
=== RUN   TestPrepareVideoEditRequest_RequiresSource
--- PASS: TestPrepareVideoEditRequest_RequiresSource (0.00s)
=== RUN   TestPrepareVideoEditRequest_ExtraParams
--- PASS: TestPrepareVideoEditRequest_ExtraParams (0.00s)
=== RUN   TestPrepareVideoEditRequest_MultipartUpload
--- PASS: TestPrepareVideoEditRequest_MultipartUpload (0.00s)
=== RUN   TestPrepareVideoEditRequest_MultipartBracketedVideoID
--- PASS: TestPrepareVideoEditRequest_MultipartBracketedVideoID (0.00s)
=== RUN   TestVideoRouteShapesDoNotConflict
--- PASS: TestVideoRouteShapesDoNotConflict (0.00s)
=== RUN   TestIsLocalhost
--- PASS: TestIsLocalhost (0.00s)
=== RUN   TestIsLocalhostOrigin
--- PASS: TestIsLocalhostOrigin (0.00s)
=== RUN   TestLoopbackRedirectURIsIPv6
--- PASS: TestLoopbackRedirectURIsIPv6 (0.00s)
=== RUN   TestPrivateUseRedirectSchemes
--- PASS: TestPrivateUseRedirectSchemes (0.00s)
=== RUN   TestMatchRedirectURIPrivateUseSchemes
--- PASS: TestMatchRedirectURIPrivateUseSchemes (0.00s)
=== RUN   TestShouldUseFilterDataCacheAllowsUnscopedEmptyQuery
--- PASS: TestShouldUseFilterDataCacheAllowsUnscopedEmptyQuery (0.00s)
=== RUN   TestParseComplexityFilters
=== RUN   TestParseComplexityFilters/parses_tier_and_mechanism
--- PASS: TestParseComplexityFilters/parses_tier_and_mechanism (0.00s)
=== RUN   TestParseComplexityFilters/leaves_filters_unchanged_when_parameters_are_absent
--- PASS: TestParseComplexityFilters/leaves_filters_unchanged_when_parameters_are_absent (0.00s)
--- PASS: TestParseComplexityFilters (0.00s)
=== RUN   TestParseToolCallNamesFilter
=== RUN   TestParseToolCallNamesFilter/parses_comma-separated_names
--- PASS: TestParseToolCallNamesFilter/parses_comma-separated_names (0.00s)
=== RUN   TestParseToolCallNamesFilter/leaves_filters_unchanged_when_parameter_is_absent
--- PASS: TestParseToolCallNamesFilter/leaves_filters_unchanged_when_parameter_is_absent (0.00s)
=== RUN   TestParseToolCallNamesFilter/histogram_filters_honour_it
--- PASS: TestParseToolCallNamesFilter/histogram_filters_honour_it (0.00s)
--- PASS: TestParseToolCallNamesFilter (0.00s)
=== RUN   TestParseParentRequestIDFilter
--- PASS: TestParseParentRequestIDFilter (0.00s)
=== RUN   TestShouldUseFilterDataCacheRejectsSearchQuery
--- PASS: TestShouldUseFilterDataCacheRejectsSearchQuery (0.00s)
=== RUN   TestShouldUseFilterDataCacheRejectsScopedContext
--- PASS: TestShouldUseFilterDataCacheRejectsScopedContext (0.00s)
=== RUN   TestGetMCPLogByIDRedactionMapping
=== RUN   TestGetMCPLogByIDRedactionMapping/no_resolver
--- PASS: TestGetMCPLogByIDRedactionMapping/no_resolver (0.01s)
=== RUN   TestGetMCPLogByIDRedactionMapping/authorized_mapping
--- PASS: TestGetMCPLogByIDRedactionMapping/authorized_mapping (0.00s)
=== RUN   TestGetMCPLogByIDRedactionMapping/resolver_error
--- PASS: TestGetMCPLogByIDRedactionMapping/resolver_error (0.00s)
--- PASS: TestGetMCPLogByIDRedactionMapping (0.01s)
=== RUN   TestShouldCacheFilterDimensions_NarrowsToRawScans
=== RUN   TestShouldCacheFilterDimensions_NarrowsToRawScans/single_matview_dimension
--- PASS: TestShouldCacheFilterDimensions_NarrowsToRawScans/single_matview_dimension (0.00s)
=== RUN   TestShouldCacheFilterDimensions_NarrowsToRawScans/several_matview_dimensions
--- PASS: TestShouldCacheFilterDimensions_NarrowsToRawScans/several_matview_dimensions (0.00s)
=== RUN   TestShouldCacheFilterDimensions_NarrowsToRawScans/metadata_keys_alone
--- PASS: TestShouldCacheFilterDimensions_NarrowsToRawScans/metadata_keys_alone (0.00s)
=== RUN   TestShouldCacheFilterDimensions_NarrowsToRawScans/metadata_keys_mixed_in
--- PASS: TestShouldCacheFilterDimensions_NarrowsToRawScans/metadata_keys_mixed_in (0.00s)
=== RUN   TestShouldCacheFilterDimensions_NarrowsToRawScans/default_all_dimensions
--- PASS: TestShouldCacheFilterDimensions_NarrowsToRawScans/default_all_dimensions (0.00s)
--- PASS: TestShouldCacheFilterDimensions_NarrowsToRawScans (0.00s)
=== RUN   TestShouldCacheFilterDimensions_NonPostgresCachesEverything
--- PASS: TestShouldCacheFilterDimensions_NonPostgresCachesEverything (0.00s)
=== RUN   TestFilterDataCacheIdentity_PartitionsPerCaller
--- PASS: TestFilterDataCacheIdentity_PartitionsPerCaller (0.00s)
=== RUN   TestGetDashboard
=== RUN   TestGetDashboard/success_includes_all_sections
--- PASS: TestGetDashboard/success_includes_all_sections (0.02s)
=== RUN   TestGetDashboard/sub-query_error_returns_no_partial_dashboard
--- PASS: TestGetDashboard/sub-query_error_returns_no_partial_dashboard (0.00s)
=== RUN   TestGetDashboard/MCP_filters_are_isolated_from_LLM_filters
--- PASS: TestGetDashboard/MCP_filters_are_isolated_from_LLM_filters (0.00s)
--- PASS: TestGetDashboard (0.02s)
=== RUN   TestRecalculateLogCostsResolvesPeriodFilter
--- PASS: TestRecalculateLogCostsResolvesPeriodFilter (0.03s)
=== RUN   TestRecalculateLogCostsRejectsDuplicateJob
--- PASS: TestRecalculateLogCostsRejectsDuplicateJob (0.00s)
=== RUN   TestCancelRecalculateCost
=== RUN   TestCancelRecalculateCost/cancels_the_in-flight_job_when_no_id_is_given
--- PASS: TestCancelRecalculateCost/cancels_the_in-flight_job_when_no_id_is_given (0.00s)
=== RUN   TestCancelRecalculateCost/cancels_the_job_named_by_id
--- PASS: TestCancelRecalculateCost/cancels_the_job_named_by_id (0.00s)
=== RUN   TestCancelRecalculateCost/an_already-terminal_job_is_returned_unchanged
--- PASS: TestCancelRecalculateCost/an_already-terminal_job_is_returned_unchanged (0.00s)
=== RUN   TestCancelRecalculateCost/refuses_to_cancel_a_job_of_another_kind
--- PASS: TestCancelRecalculateCost/refuses_to_cancel_a_job_of_another_kind (0.00s)
=== RUN   TestCancelRecalculateCost/404_when_there_is_nothing_to_cancel
--- PASS: TestCancelRecalculateCost/404_when_there_is_nothing_to_cancel (0.00s)
=== RUN   TestCancelRecalculateCost/503_when_the_background_runner_is_not_wired
--- PASS: TestCancelRecalculateCost/503_when_the_background_runner_is_not_wired (0.00s)
--- PASS: TestCancelRecalculateCost (0.00s)
=== RUN   TestProjectFilterReachesEveryLogQuery
=== RUN   TestProjectFilterReachesEveryLogQuery/list
--- PASS: TestProjectFilterReachesEveryLogQuery/list (0.00s)
=== RUN   TestProjectFilterReachesEveryLogQuery/stats
--- PASS: TestProjectFilterReachesEveryLogQuery/stats (0.00s)
=== RUN   TestProjectFilterReachesEveryLogQuery/histograms
--- PASS: TestProjectFilterReachesEveryLogQuery/histograms (0.00s)
--- PASS: TestProjectFilterReachesEveryLogQuery (0.00s)
=== RUN   TestFilterDataListsProjects
--- PASS: TestFilterDataListsProjects (0.00s)
=== RUN   TestMCPAttributionFilterParsing
--- PASS: TestMCPAttributionFilterParsing (0.00s)
=== RUN   TestLogsStatsWithoutCompareIsUnchanged
--- PASS: TestLogsStatsWithoutCompareIsUnchanged (0.00s)
=== RUN   TestLogsStatsComparePreviousWindow
--- PASS: TestLogsStatsComparePreviousWindow (0.00s)
=== RUN   TestLogsStatsCompareSkippedWhenUnbounded
--- PASS: TestLogsStatsCompareSkippedWhenUnbounded (0.00s)
=== RUN   TestLogsStatsCompareDegradesOnPreviousError
--- PASS: TestLogsStatsCompareDegradesOnPreviousError (0.00s)
=== RUN   TestLogsStatsCompareSkippedForRequestIDLookup
--- PASS: TestLogsStatsCompareSkippedForRequestIDLookup (0.00s)
=== RUN   TestLogsStatsCompareSkippedForEmptyWindow
=== RUN   TestLogsStatsCompareSkippedForEmptyWindow/zero_length
--- PASS: TestLogsStatsCompareSkippedForEmptyWindow/zero_length (0.00s)
=== RUN   TestLogsStatsCompareSkippedForEmptyWindow/reversed
--- PASS: TestLogsStatsCompareSkippedForEmptyWindow/reversed (0.00s)
--- PASS: TestLogsStatsCompareSkippedForEmptyWindow (0.00s)
=== RUN   TestHiddenRequestTypesRoutes
--- PASS: TestHiddenRequestTypesRoutes (0.00s)
=== RUN   TestHiddenRequestTypesFilterCache
--- PASS: TestHiddenRequestTypesFilterCache (0.00s)
=== RUN   TestMCPClientRequestReadsAllowByDefaultUnderBothKeys
=== RUN   TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/current_key_only
--- PASS: TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/current_key_only (0.00s)
=== RUN   TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/earlier_key_only
--- PASS: TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/earlier_key_only (0.00s)
=== RUN   TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/both_sent,_current_wins
--- PASS: TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/both_sent,_current_wins (0.00s)
=== RUN   TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/neither_sent
--- PASS: TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/neither_sent (0.00s)
=== RUN   TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/fields_of_the_outer_request_still_decode
--- PASS: TestMCPClientRequestReadsAllowByDefaultUnderBothKeys/fields_of_the_outer_request_still_decode (0.00s)
--- PASS: TestMCPClientRequestReadsAllowByDefaultUnderBothKeys (0.00s)
=== RUN   TestMCPClientUpdateRequestResolvedAllowByDefault
=== RUN   TestMCPClientUpdateRequestResolvedAllowByDefault/neither_sent_keeps_existing_true
--- PASS: TestMCPClientUpdateRequestResolvedAllowByDefault/neither_sent_keeps_existing_true (0.00s)
=== RUN   TestMCPClientUpdateRequestResolvedAllowByDefault/neither_sent_keeps_existing_false
--- PASS: TestMCPClientUpdateRequestResolvedAllowByDefault/neither_sent_keeps_existing_false (0.00s)
=== RUN   TestMCPClientUpdateRequestResolvedAllowByDefault/current_key_only
--- PASS: TestMCPClientUpdateRequestResolvedAllowByDefault/current_key_only (0.00s)
=== RUN   TestMCPClientUpdateRequestResolvedAllowByDefault/earlier_key_only
--- PASS: TestMCPClientUpdateRequestResolvedAllowByDefault/earlier_key_only (0.00s)
=== RUN   TestMCPClientUpdateRequestResolvedAllowByDefault/both_sent,_current_wins
--- PASS: TestMCPClientUpdateRequestResolvedAllowByDefault/both_sent,_current_wins (0.00s)
--- PASS: TestMCPClientUpdateRequestResolvedAllowByDefault (0.00s)
=== RUN   TestParseOauthScopesJSON
=== RUN   TestParseOauthScopesJSON/empty_string
--- PASS: TestParseOauthScopesJSON/empty_string (0.00s)
=== RUN   TestParseOauthScopesJSON/whitespace_only
--- PASS: TestParseOauthScopesJSON/whitespace_only (0.00s)
=== RUN   TestParseOauthScopesJSON/json_null_(provider_omitted_scope)
--- PASS: TestParseOauthScopesJSON/json_null_(provider_omitted_scope) (0.00s)
=== RUN   TestParseOauthScopesJSON/empty_array
--- PASS: TestParseOauthScopesJSON/empty_array (0.00s)
=== RUN   TestParseOauthScopesJSON/malformed
--- PASS: TestParseOauthScopesJSON/malformed (0.00s)
=== RUN   TestParseOauthScopesJSON/wrong_shape
--- PASS: TestParseOauthScopesJSON/wrong_shape (0.00s)
=== RUN   TestParseOauthScopesJSON/list
--- PASS: TestParseOauthScopesJSON/list (0.00s)
=== RUN   TestParseOauthScopesJSON/drops_blank_entries
--- PASS: TestParseOauthScopesJSON/drops_blank_entries (0.00s)
=== RUN   TestParseOauthScopesJSON/trims_entries
--- PASS: TestParseOauthScopesJSON/trims_entries (0.00s)
--- PASS: TestParseOauthScopesJSON (0.00s)
=== RUN   TestTokenRow_TokenOnlyFields
--- PASS: TestTokenRow_TokenOnlyFields (0.00s)
=== RUN   TestMCPClientCredentialResponse
=== RUN   TestMCPClientCredentialResponse/oauth_reads_the_shared_token,_not_the_admin_one
--- PASS: TestMCPClientCredentialResponse/oauth_reads_the_shared_token,_not_the_admin_one (0.00s)
=== RUN   TestMCPClientCredentialResponse/per_user_oauth_reads_the_admin_token,_not_the_shared_one
--- PASS: TestMCPClientCredentialResponse/per_user_oauth_reads_the_admin_token,_not_the_shared_one (0.00s)
=== RUN   TestMCPClientCredentialResponse/token_exchange_reads_the_admin_token
--- PASS: TestMCPClientCredentialResponse/token_exchange_reads_the_admin_token (0.00s)
=== RUN   TestMCPClientCredentialResponse/per_user_headers_reads_the_admin_header_credential_with_sorted_key_names
--- PASS: TestMCPClientCredentialResponse/per_user_headers_reads_the_admin_header_credential_with_sorted_key_names (0.00s)
=== RUN   TestMCPClientCredentialResponse/none_has_no_self-held_credential
--- PASS: TestMCPClientCredentialResponse/none_has_no_self-held_credential (0.00s)
=== RUN   TestMCPClientCredentialResponse/headers_has_no_self-held_credential
--- PASS: TestMCPClientCredentialResponse/headers_has_no_self-held_credential (0.00s)
=== RUN   TestMCPClientCredentialResponse/missing_rows_yield_no_block
--- PASS: TestMCPClientCredentialResponse/missing_rows_yield_no_block (0.00s)
=== RUN   TestMCPClientCredentialResponse/unreadable_header_values_still_project_status_and_timestamps
--- PASS: TestMCPClientCredentialResponse/unreadable_header_values_still_project_status_and_timestamps (0.00s)
--- PASS: TestMCPClientCredentialResponse (0.00s)
=== RUN   TestUpdateMCPClient_DisabledToEnabled_WithInvalidReplacementHeaders_PreflightRejects
--- PASS: TestUpdateMCPClient_DisabledToEnabled_WithInvalidReplacementHeaders_PreflightRejects (0.00s)
=== RUN   TestEndpointSlugDerivable
--- PASS: TestEndpointSlugDerivable (0.00s)
=== RUN   TestMCPHeadersEqual
=== RUN   TestMCPHeadersEqual/both_empty
--- PASS: TestMCPHeadersEqual/both_empty (0.00s)
=== RUN   TestMCPHeadersEqual/identical_single_header
--- PASS: TestMCPHeadersEqual/identical_single_header (0.00s)
=== RUN   TestMCPHeadersEqual/changed_value
--- PASS: TestMCPHeadersEqual/changed_value (0.00s)
=== RUN   TestMCPHeadersEqual/added_header
--- PASS: TestMCPHeadersEqual/added_header (0.00s)
=== RUN   TestMCPHeadersEqual/removed_header
--- PASS: TestMCPHeadersEqual/removed_header (0.00s)
=== RUN   TestMCPHeadersEqual/renamed_key_with_same_value
--- PASS: TestMCPHeadersEqual/renamed_key_with_same_value (0.00s)
=== RUN   TestMCPHeadersEqual/case-differing_key
--- PASS: TestMCPHeadersEqual/case-differing_key (0.00s)
=== RUN   TestMCPHeadersEqual/empty_vs_populated
--- PASS: TestMCPHeadersEqual/empty_vs_populated (0.00s)
=== RUN   TestMCPHeadersEqual/nil_vs_empty_are_both_no_headers
--- PASS: TestMCPHeadersEqual/nil_vs_empty_are_both_no_headers (0.00s)
--- PASS: TestMCPHeadersEqual (0.00s)
=== RUN   TestIsPrematureOAuthCompletion
=== RUN   TestIsPrematureOAuthCompletion/nil_flow_(already_cleaned_up_by_a_genuine_completion)_is_not_premature
--- PASS: TestIsPrematureOAuthCompletion/nil_flow_(already_cleaned_up_by_a_genuine_completion)_is_not_premature (0.00s)
=== RUN   TestIsPrematureOAuthCompletion/pending,_unexpired_flow_means_the_callback_never_ran:_premature
--- PASS: TestIsPrematureOAuthCompletion/pending,_unexpired_flow_means_the_callback_never_ran:_premature (0.00s)
=== RUN   TestIsPrematureOAuthCompletion/pending_but_expired_flow_(abandoned_attempt)_does_not_block_a_later_fresh_completion
--- PASS: TestIsPrematureOAuthCompletion/pending_but_expired_flow_(abandoned_attempt)_does_not_block_a_later_fresh_completion (0.00s)
=== RUN   TestIsPrematureOAuthCompletion/failed_flow_row_is_not_premature_(already_resolved,_just_unsuccessfully)
--- PASS: TestIsPrematureOAuthCompletion/failed_flow_row_is_not_premature_(already_resolved,_just_unsuccessfully) (0.00s)
--- PASS: TestIsPrematureOAuthCompletion (0.00s)
=== RUN   TestProjectMCPCredentialState
=== RUN   TestProjectMCPCredentialState/per_user_oauth_connected_with_dead_admin_token_projects_needs_reauth
--- PASS: TestProjectMCPCredentialState/per_user_oauth_connected_with_dead_admin_token_projects_needs_reauth (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_oauth_connected_with_active_admin_token_stays_connected
--- PASS: TestProjectMCPCredentialState/per_user_oauth_connected_with_active_admin_token_stays_connected (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_oauth_connected_with_no_admin_token_stays_connected_(pre-retention_client)
--- PASS: TestProjectMCPCredentialState/per_user_oauth_connected_with_no_admin_token_stays_connected_(pre-retention_client) (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_oauth_connected_with_orphaned_admin_token_stays_connected
--- PASS: TestProjectMCPCredentialState/per_user_oauth_connected_with_orphaned_admin_token_stays_connected (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_oauth_ignores_the_header_credential_status
--- PASS: TestProjectMCPCredentialState/per_user_oauth_ignores_the_header_credential_status (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_headers_connected_with_stale_admin_credential_projects_needs_reauth
--- PASS: TestProjectMCPCredentialState/per_user_headers_connected_with_stale_admin_credential_projects_needs_reauth (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_headers_connected_with_active_admin_credential_stays_connected
--- PASS: TestProjectMCPCredentialState/per_user_headers_connected_with_active_admin_credential_stays_connected (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_headers_connected_with_no_admin_credential_stays_connected_(pre-retention_client)
--- PASS: TestProjectMCPCredentialState/per_user_headers_connected_with_no_admin_credential_stays_connected_(pre-retention_client) (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_headers_ignores_the_oauth_token_status
--- PASS: TestProjectMCPCredentialState/per_user_headers_ignores_the_oauth_token_status (0.00s)
=== RUN   TestProjectMCPCredentialState/pending_verification_runtime_state_passes_through_even_with_a_dead_admin_token
--- PASS: TestProjectMCPCredentialState/pending_verification_runtime_state_passes_through_even_with_a_dead_admin_token (0.00s)
=== RUN   TestProjectMCPCredentialState/disabled_runtime_state_passes_through_even_with_a_stale_admin_credential
--- PASS: TestProjectMCPCredentialState/disabled_runtime_state_passes_through_even_with_a_stale_admin_credential (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_oauth_unstable_with_a_dead_admin_token_projects_needs_reauth_(needs_reauth_is_a_bigger_indicator_than_unstable)
--- PASS: TestProjectMCPCredentialState/per_user_oauth_unstable_with_a_dead_admin_token_projects_needs_reauth_(needs_reauth_is_a_bigger_indicator_than_unstable) (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_headers_unstable_with_a_stale_admin_credential_projects_needs_reauth
--- PASS: TestProjectMCPCredentialState/per_user_headers_unstable_with_a_stale_admin_credential_projects_needs_reauth (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_oauth_unstable_with_an_active_admin_token_stays_unstable
--- PASS: TestProjectMCPCredentialState/per_user_oauth_unstable_with_an_active_admin_token_stays_unstable (0.00s)
=== RUN   TestProjectMCPCredentialState/shared_oauth_connected_with_rotated_token_projects_needs_reauth
--- PASS: TestProjectMCPCredentialState/shared_oauth_connected_with_rotated_token_projects_needs_reauth (0.00s)
=== RUN   TestProjectMCPCredentialState/shared_oauth_connected_with_active_token_stays_connected
--- PASS: TestProjectMCPCredentialState/shared_oauth_connected_with_active_token_stays_connected (0.00s)
=== RUN   TestProjectMCPCredentialState/shared_oauth_connected_with_no_token_row_stays_connected
--- PASS: TestProjectMCPCredentialState/shared_oauth_connected_with_no_token_row_stays_connected (0.00s)
=== RUN   TestProjectMCPCredentialState/shared_oauth_unstable_with_rotated_token_projects_needs_reauth
--- PASS: TestProjectMCPCredentialState/shared_oauth_unstable_with_rotated_token_projects_needs_reauth (0.00s)
=== RUN   TestProjectMCPCredentialState/shared_oauth_ignores_the_per-user_admin_credential_statuses
--- PASS: TestProjectMCPCredentialState/shared_oauth_ignores_the_per-user_admin_credential_statuses (0.00s)
=== RUN   TestProjectMCPCredentialState/per_user_oauth_ignores_the_shared_token_status
--- PASS: TestProjectMCPCredentialState/per_user_oauth_ignores_the_shared_token_status (0.00s)
=== RUN   TestProjectMCPCredentialState/shared_headers_are_never_projected
--- PASS: TestProjectMCPCredentialState/shared_headers_are_never_projected (0.00s)
=== RUN   TestProjectMCPCredentialState/disabled_shared_oauth_passes_through_even_with_a_rotated_token
--- PASS: TestProjectMCPCredentialState/disabled_shared_oauth_passes_through_even_with_a_rotated_token (0.00s)
=== RUN   TestProjectMCPCredentialState/non-auth_clients_are_never_projected
--- PASS: TestProjectMCPCredentialState/non-auth_clients_are_never_projected (0.00s)
--- PASS: TestProjectMCPCredentialState (0.00s)
=== RUN   TestToolSyncIntervalBoundsRejectNanosecondValues
--- PASS: TestToolSyncIntervalBoundsRejectNanosecondValues (0.00s)
=== RUN   TestToolSyncIntervalBoundsRejectNegatives
--- PASS: TestToolSyncIntervalBoundsRejectNegatives (0.00s)
=== RUN   TestToolSyncIntervalBoundsAcceptRealisticValues
--- PASS: TestToolSyncIntervalBoundsAcceptRealisticValues (0.00s)
=== RUN   TestToolSyncIntervalBoundsSurvivePersistence
--- PASS: TestToolSyncIntervalBoundsSurvivePersistence (0.00s)
=== RUN   TestToolSyncIntervalBoundsPreventOverflow
--- PASS: TestToolSyncIntervalBoundsPreventOverflow (0.00s)
=== RUN   TestUpdateMCPClientCredentialsWithRetry_NotApplicable_ReturnsImmediately
--- PASS: TestUpdateMCPClientCredentialsWithRetry_NotApplicable_ReturnsImmediately (0.00s)
=== RUN   TestUpdateMCPClientCredentialsWithRetry_TransientReconnectError_StillRetries
--- PASS: TestUpdateMCPClientCredentialsWithRetry_TransientReconnectError_StillRetries (1.00s)
=== RUN   TestIsConcurrentVerifyReplay
=== RUN   TestIsConcurrentVerifyReplay/genuine_concurrent_race:_started_pending_verification,_another_request_finished_the_bootstrap_first
--- PASS: TestIsConcurrentVerifyReplay/genuine_concurrent_race:_started_pending_verification,_another_request_finished_the_bootstrap_first (0.00s)
=== RUN   TestIsConcurrentVerifyReplay/repair_on_an_already-verified_client:_not_pending_verification,_tools_already_present
--- PASS: TestIsConcurrentVerifyReplay/repair_on_an_already-verified_client:_not_pending_verification,_tools_already_present (0.00s)
=== RUN   TestIsConcurrentVerifyReplay/first-time_verification_with_no_race:_started_pending_verification,_still_no_tools_after_reload
--- PASS: TestIsConcurrentVerifyReplay/first-time_verification_with_no_race:_started_pending_verification,_still_no_tools_after_reload (0.00s)
--- PASS: TestIsConcurrentVerifyReplay (0.00s)
=== RUN   TestConsentFlowDetail
=== RUN   TestConsentFlowDetail/pending_flow_returns_client_and_modes
--- PASS: TestConsentFlowDetail/pending_flow_returns_client_and_modes (0.00s)
=== RUN   TestConsentFlowDetail/missing_flow_returns_404
--- PASS: TestConsentFlowDetail/missing_flow_returns_404 (0.00s)
=== RUN   TestConsentFlowDetail/empty_flow_id_returns_400
--- PASS: TestConsentFlowDetail/empty_flow_id_returns_400 (0.00s)
=== RUN   TestConsentFlowDetail/expired_flow_returns_410
--- PASS: TestConsentFlowDetail/expired_flow_returns_410 (0.00s)
=== RUN   TestConsentFlowDetail/already-consented_flow_returns_410
--- PASS: TestConsentFlowDetail/already-consented_flow_returns_410 (0.00s)
--- PASS: TestConsentFlowDetail (0.00s)
=== RUN   TestConsentAvailableModes
=== RUN   TestConsentAvailableModes/vk_and_session_when_auth_not_enforced
--- PASS: TestConsentAvailableModes/vk_and_session_when_auth_not_enforced (0.00s)
=== RUN   TestConsentAvailableModes/vk_only_when_auth_enforced
--- PASS: TestConsentAvailableModes/vk_only_when_auth_enforced (0.00s)
=== RUN   TestConsentAvailableModes/adds_user_when_resolver_offers_it
--- PASS: TestConsentAvailableModes/adds_user_when_resolver_offers_it (0.00s)
=== RUN   TestConsentAvailableModes/disable_vk_drops_vk_when_user_mode_available
--- PASS: TestConsentAvailableModes/disable_vk_drops_vk_when_user_mode_available (0.00s)
=== RUN   TestConsentAvailableModes/disable_vk_leaves_user-only_when_auth_enforced
--- PASS: TestConsentAvailableModes/disable_vk_leaves_user-only_when_auth_enforced (0.00s)
=== RUN   TestConsentAvailableModes/disable_vk_ignored_without_user_mode
--- PASS: TestConsentAvailableModes/disable_vk_ignored_without_user_mode (0.00s)
--- PASS: TestConsentAvailableModes (0.00s)
=== RUN   TestConsentFlowSubmit_VK
=== RUN   TestConsentFlowSubmit_VK/active_VK_mints_a_code
--- PASS: TestConsentFlowSubmit_VK/active_VK_mints_a_code (0.00s)
=== RUN   TestConsentFlowSubmit_VK/inactive_VK_is_rejected
--- PASS: TestConsentFlowSubmit_VK/inactive_VK_is_rejected (0.00s)
=== RUN   TestConsentFlowSubmit_VK/unknown_VK_is_rejected
--- PASS: TestConsentFlowSubmit_VK/unknown_VK_is_rejected (0.00s)
=== RUN   TestConsentFlowSubmit_VK/empty_VK_value_is_rejected
--- PASS: TestConsentFlowSubmit_VK/empty_VK_value_is_rejected (0.00s)
=== RUN   TestConsentFlowSubmit_VK/double_submit_returns_410_on_the_second_attempt
--- PASS: TestConsentFlowSubmit_VK/double_submit_returns_410_on_the_second_attempt (0.00s)
--- PASS: TestConsentFlowSubmit_VK (0.00s)
=== RUN   TestConsentFlowSubmit_Session
=== RUN   TestConsentFlowSubmit_Session/session_mode_mints_a_server-side_token_when_auth_not_enforced
--- PASS: TestConsentFlowSubmit_Session/session_mode_mints_a_server-side_token_when_auth_not_enforced (0.00s)
=== RUN   TestConsentFlowSubmit_Session/session_mode_is_unavailable_when_auth_is_enforced
--- PASS: TestConsentFlowSubmit_Session/session_mode_is_unavailable_when_auth_is_enforced (0.00s)
--- PASS: TestConsentFlowSubmit_Session (0.00s)
=== RUN   TestConsentFlowSubmit_User
=== RUN   TestConsentFlowSubmit_User/user_mode_rejected_when_no_resolver_(mode_not_offered)
--- PASS: TestConsentFlowSubmit_User/user_mode_rejected_when_no_resolver_(mode_not_offered) (0.00s)
=== RUN   TestConsentFlowSubmit_User/resolved_session_yields_user_mode
--- PASS: TestConsentFlowSubmit_User/resolved_session_yields_user_mode (0.00s)
=== RUN   TestConsentFlowSubmit_User/user_mode_without_a_session_is_rejected
--- PASS: TestConsentFlowSubmit_User/user_mode_without_a_session_is_rejected (0.00s)
--- PASS: TestConsentFlowSubmit_User (0.00s)
=== RUN   TestConsentFlowSubmit_VKUserBinding
=== RUN   TestConsentFlowSubmit_VKUserBinding/bound_VK_upgrades_to_user_when_logged-in_user_matches
--- PASS: TestConsentFlowSubmit_VKUserBinding/bound_VK_upgrades_to_user_when_logged-in_user_matches (0.00s)
=== RUN   TestConsentFlowSubmit_VKUserBinding/bound_VK_rejected_when_logged-in_user_differs
--- PASS: TestConsentFlowSubmit_VKUserBinding/bound_VK_rejected_when_logged-in_user_differs (0.00s)
=== RUN   TestConsentFlowSubmit_VKUserBinding/bound_VK_rejected_when_not_signed_in
--- PASS: TestConsentFlowSubmit_VKUserBinding/bound_VK_rejected_when_not_signed_in (0.00s)
--- PASS: TestConsentFlowSubmit_VKUserBinding (0.00s)
=== RUN   TestDiscovery_GatedOnAuthMode
=== RUN   TestDiscovery_GatedOnAuthMode/headers_mode_returns_404_on_all_discovery_endpoints
--- PASS: TestDiscovery_GatedOnAuthMode/headers_mode_returns_404_on_all_discovery_endpoints (0.00s)
=== RUN   TestDiscovery_GatedOnAuthMode/oauth_and_both_modes_serve_discovery
--- PASS: TestDiscovery_GatedOnAuthMode/oauth_and_both_modes_serve_discovery (0.00s)
--- PASS: TestDiscovery_GatedOnAuthMode (0.07s)
=== RUN   TestDiscovery_ProtectedResourceMetadata
--- PASS: TestDiscovery_ProtectedResourceMetadata (0.00s)
=== RUN   TestDiscovery_AuthorizationServerMetadata
--- PASS: TestDiscovery_AuthorizationServerMetadata (0.00s)
=== RUN   TestDiscovery_JWKS
--- PASS: TestDiscovery_JWKS (0.09s)
=== RUN   TestHandleRegister_DCR
=== RUN   TestHandleRegister_DCR/valid_registration_returns_201_with_defaults
    mcpoauth2issuance_test.go:123: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:123
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleRegister_DCR/valid_registration_returns_201_with_defaults
--- FAIL: TestHandleRegister_DCR/valid_registration_returns_201_with_defaults (0.00s)
=== RUN   TestHandleRegister_DCR/missing_redirect_uris_is_rejected
    mcpoauth2issuance_test.go:140: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:140
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleRegister_DCR/missing_redirect_uris_is_rejected
--- FAIL: TestHandleRegister_DCR/missing_redirect_uris_is_rejected (0.00s)
=== RUN   TestHandleRegister_DCR/non-public_auth_method_is_rejected
    mcpoauth2issuance_test.go:151: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:151
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleRegister_DCR/non-public_auth_method_is_rejected
--- FAIL: TestHandleRegister_DCR/non-public_auth_method_is_rejected (0.00s)
=== RUN   TestHandleRegister_DCR/malformed_JSON_is_rejected
    mcpoauth2issuance_test.go:162: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:162
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleRegister_DCR/malformed_JSON_is_rejected
--- FAIL: TestHandleRegister_DCR/malformed_JSON_is_rejected (0.00s)
=== RUN   TestHandleRegister_DCR/unsupported_grant_type_is_rejected
    mcpoauth2issuance_test.go:173: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:173
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleRegister_DCR/unsupported_grant_type_is_rejected
--- FAIL: TestHandleRegister_DCR/unsupported_grant_type_is_rejected (0.00s)
=== RUN   TestHandleRegister_DCR/unsupported_response_type_is_rejected
    mcpoauth2issuance_test.go:184: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:184
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleRegister_DCR/unsupported_response_type_is_rejected
--- FAIL: TestHandleRegister_DCR/unsupported_response_type_is_rejected (0.00s)
--- FAIL: TestHandleRegister_DCR (0.01s)
=== RUN   TestHandleAuthorize
=== RUN   TestHandleAuthorize/happy_path_redirects_to_the_consent_page
    mcpoauth2issuance_test.go:209: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:209
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/happy_path_redirects_to_the_consent_page
--- FAIL: TestHandleAuthorize/happy_path_redirects_to_the_consent_page (0.00s)
=== RUN   TestHandleAuthorize/loopback_redirect_matches_on_any_port
    mcpoauth2issuance_test.go:219: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:219
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/loopback_redirect_matches_on_any_port
--- FAIL: TestHandleAuthorize/loopback_redirect_matches_on_any_port (0.00s)
=== RUN   TestHandleAuthorize/unknown_client_id_is_rejected
    mcpoauth2issuance_test.go:229: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:229
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/unknown_client_id_is_rejected
--- FAIL: TestHandleAuthorize/unknown_client_id_is_rejected (0.00s)
=== RUN   TestHandleAuthorize/unregistered_redirect_uri_is_rejected
    mcpoauth2issuance_test.go:238: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:238
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/unregistered_redirect_uri_is_rejected
--- FAIL: TestHandleAuthorize/unregistered_redirect_uri_is_rejected (0.00s)
=== RUN   TestHandleAuthorize/non-code_response_type_redirects_with_error
    mcpoauth2issuance_test.go:260: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:260
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/non-code_response_type_redirects_with_error
--- FAIL: TestHandleAuthorize/non-code_response_type_redirects_with_error (0.00s)
=== RUN   TestHandleAuthorize/non-S256_challenge_method_redirects_with_error
    mcpoauth2issuance_test.go:260: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:260
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/non-S256_challenge_method_redirects_with_error
--- FAIL: TestHandleAuthorize/non-S256_challenge_method_redirects_with_error (0.00s)
=== RUN   TestHandleAuthorize/mismatched_resource_redirects_with_error
    mcpoauth2issuance_test.go:260: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:260
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/mismatched_resource_redirects_with_error
--- FAIL: TestHandleAuthorize/mismatched_resource_redirects_with_error (0.00s)
=== RUN   TestHandleAuthorize/scope_exceeds_registered_redirects_with_error
    mcpoauth2issuance_test.go:260: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:260
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/scope_exceeds_registered_redirects_with_error
--- FAIL: TestHandleAuthorize/scope_exceeds_registered_redirects_with_error (0.00s)
=== RUN   TestHandleAuthorize/omitted_resource_defaults_to_canonical_and_proceeds
    mcpoauth2issuance_test.go:275: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:275
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize/omitted_resource_defaults_to_canonical_and_proceeds
--- FAIL: TestHandleAuthorize/omitted_resource_defaults_to_canonical_and_proceeds (0.00s)
--- FAIL: TestHandleAuthorize (0.01s)
=== RUN   TestHandleToken_AuthorizationCode
=== RUN   TestHandleToken_AuthorizationCode/happy_path_issues_a_verifiable_token_pair
    mcpoauth2issuance_test.go:292: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:292
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_AuthorizationCode/happy_path_issues_a_verifiable_token_pair
--- FAIL: TestHandleToken_AuthorizationCode/happy_path_issues_a_verifiable_token_pair (0.00s)
=== RUN   TestHandleToken_AuthorizationCode/PKCE_mismatch_is_rejected
    mcpoauth2issuance_test.go:322: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:322
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_AuthorizationCode/PKCE_mismatch_is_rejected
--- FAIL: TestHandleToken_AuthorizationCode/PKCE_mismatch_is_rejected (0.00s)
=== RUN   TestHandleToken_AuthorizationCode/code_is_single-use
    mcpoauth2issuance_test.go:339: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:339
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_AuthorizationCode/code_is_single-use
--- FAIL: TestHandleToken_AuthorizationCode/code_is_single-use (0.00s)
=== RUN   TestHandleToken_AuthorizationCode/expired_code_is_rejected
    mcpoauth2issuance_test.go:361: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:361
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_AuthorizationCode/expired_code_is_rejected
--- FAIL: TestHandleToken_AuthorizationCode/expired_code_is_rejected (0.00s)
=== RUN   TestHandleToken_AuthorizationCode/client_id_mismatch_is_rejected
    mcpoauth2issuance_test.go:378: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:378
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_AuthorizationCode/client_id_mismatch_is_rejected
--- FAIL: TestHandleToken_AuthorizationCode/client_id_mismatch_is_rejected (0.00s)
=== RUN   TestHandleToken_AuthorizationCode/missing_redirect_uri_is_rejected
    mcpoauth2issuance_test.go:395: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:395
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_AuthorizationCode/missing_redirect_uri_is_rejected
--- FAIL: TestHandleToken_AuthorizationCode/missing_redirect_uri_is_rejected (0.00s)
=== RUN   TestHandleToken_AuthorizationCode/missing_required_fields_are_rejected
    mcpoauth2issuance_test.go:411: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:411
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_AuthorizationCode/missing_required_fields_are_rejected
--- FAIL: TestHandleToken_AuthorizationCode/missing_required_fields_are_rejected (0.00s)
=== RUN   TestHandleToken_AuthorizationCode/unsupported_grant_type_is_rejected
    mcpoauth2issuance_test.go:419: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:419
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_AuthorizationCode/unsupported_grant_type_is_rejected
--- FAIL: TestHandleToken_AuthorizationCode/unsupported_grant_type_is_rejected (0.00s)
--- FAIL: TestHandleToken_AuthorizationCode (0.01s)
=== RUN   TestHandleToken_RefreshRotationAndReplay
=== RUN   TestHandleToken_RefreshRotationAndReplay/rotation_issues_a_new_pair_and_carries_the_family
    mcpoauth2issuance_test.go:441: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:441
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_RefreshRotationAndReplay/rotation_issues_a_new_pair_and_carries_the_family
--- FAIL: TestHandleToken_RefreshRotationAndReplay/rotation_issues_a_new_pair_and_carries_the_family (0.00s)
=== RUN   TestHandleToken_RefreshRotationAndReplay/replaying_a_rotated_token_revokes_the_whole_family
    mcpoauth2issuance_test.go:468: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:468
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_RefreshRotationAndReplay/replaying_a_rotated_token_revokes_the_whole_family
--- FAIL: TestHandleToken_RefreshRotationAndReplay/replaying_a_rotated_token_revokes_the_whole_family (0.00s)
=== RUN   TestHandleToken_RefreshRotationAndReplay/client_id_mismatch_is_rejected
    mcpoauth2issuance_test.go:499: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:499
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_RefreshRotationAndReplay/client_id_mismatch_is_rejected
--- FAIL: TestHandleToken_RefreshRotationAndReplay/client_id_mismatch_is_rejected (0.00s)
=== RUN   TestHandleToken_RefreshRotationAndReplay/missing_fields_are_rejected
    mcpoauth2issuance_test.go:512: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:512
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_RefreshRotationAndReplay/missing_fields_are_rejected
--- FAIL: TestHandleToken_RefreshRotationAndReplay/missing_fields_are_rejected (0.00s)
--- FAIL: TestHandleToken_RefreshRotationAndReplay (0.01s)
=== RUN   TestHandleToken_RefreshVKIdentityDisabled
=== RUN   TestHandleToken_RefreshVKIdentityDisabled/vk_refresh_rejected_when_disabled_and_user_mode_available
    mcpoauth2issuance_test.go:538: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:538
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_RefreshVKIdentityDisabled/vk_refresh_rejected_when_disabled_and_user_mode_available
--- FAIL: TestHandleToken_RefreshVKIdentityDisabled/vk_refresh_rejected_when_disabled_and_user_mode_available (0.00s)
=== RUN   TestHandleToken_RefreshVKIdentityDisabled/vk_refresh_not_blocked_by_flag_when_user_mode_unavailable
    mcpoauth2issuance_test.go:554: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:554
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_RefreshVKIdentityDisabled/vk_refresh_not_blocked_by_flag_when_user_mode_unavailable
--- FAIL: TestHandleToken_RefreshVKIdentityDisabled/vk_refresh_not_blocked_by_flag_when_user_mode_unavailable (0.00s)
--- FAIL: TestHandleToken_RefreshVKIdentityDisabled (0.00s)
=== RUN   TestHandleToken_RefreshUserLiveness
=== RUN   TestHandleToken_RefreshUserLiveness/user_refresh_rejected_when_the_user_is_no_longer_active
    mcpoauth2issuance_test.go:591: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:591
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_RefreshUserLiveness/user_refresh_rejected_when_the_user_is_no_longer_active
--- FAIL: TestHandleToken_RefreshUserLiveness/user_refresh_rejected_when_the_user_is_no_longer_active (0.00s)
=== RUN   TestHandleToken_RefreshUserLiveness/user_refresh_succeeds_when_the_user_is_active
    mcpoauth2issuance_test.go:608: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:608
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleToken_RefreshUserLiveness/user_refresh_succeeds_when_the_user_is_active
--- FAIL: TestHandleToken_RefreshUserLiveness/user_refresh_succeeds_when_the_user_is_active (0.00s)
--- FAIL: TestHandleToken_RefreshUserLiveness (0.00s)
=== RUN   TestUpdateConfig_RejectsAuthCodeTTLAboveMax
=== RUN   TestUpdateConfig_RejectsAuthCodeTTLAboveMax/oauth
    mcpoauth2issuance_test.go:644: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:644
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestUpdateConfig_RejectsAuthCodeTTLAboveMax/oauth
--- FAIL: TestUpdateConfig_RejectsAuthCodeTTLAboveMax/oauth (0.00s)
=== RUN   TestUpdateConfig_RejectsAuthCodeTTLAboveMax/both
    mcpoauth2issuance_test.go:644: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:644
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestUpdateConfig_RejectsAuthCodeTTLAboveMax/both
--- FAIL: TestUpdateConfig_RejectsAuthCodeTTLAboveMax/both (0.00s)
=== RUN   TestUpdateConfig_RejectsAuthCodeTTLAboveMax/headers
    mcpoauth2issuance_test.go:644: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:644
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestUpdateConfig_RejectsAuthCodeTTLAboveMax/headers
--- FAIL: TestUpdateConfig_RejectsAuthCodeTTLAboveMax/headers (0.00s)
--- FAIL: TestUpdateConfig_RejectsAuthCodeTTLAboveMax (0.01s)
=== RUN   TestHandleAuthorize_AuthCodeTTLResolution
=== RUN   TestHandleAuthorize_AuthCodeTTLResolution/above_cap_is_clamped
    mcpoauth2issuance_test.go:676: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:676
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize_AuthCodeTTLResolution/above_cap_is_clamped
--- FAIL: TestHandleAuthorize_AuthCodeTTLResolution/above_cap_is_clamped (0.00s)
=== RUN   TestHandleAuthorize_AuthCodeTTLResolution/zero_falls_back_to_default
    mcpoauth2issuance_test.go:676: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:676
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize_AuthCodeTTLResolution/zero_falls_back_to_default
--- FAIL: TestHandleAuthorize_AuthCodeTTLResolution/zero_falls_back_to_default (0.00s)
=== RUN   TestHandleAuthorize_AuthCodeTTLResolution/in-range_used_verbatim
    mcpoauth2issuance_test.go:676: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:31
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:39
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/mcpoauth2issuance_test.go:676
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestHandleAuthorize_AuthCodeTTLResolution/in-range_used_verbatim
--- FAIL: TestHandleAuthorize_AuthCodeTTLResolution/in-range_used_verbatim (0.00s)
--- FAIL: TestHandleAuthorize_AuthCodeTTLResolution (0.00s)
=== RUN   TestExtractBearerJWT
=== RUN   TestExtractBearerJWT/jwt_bearer
--- PASS: TestExtractBearerJWT/jwt_bearer (0.00s)
=== RUN   TestExtractBearerJWT/case-insensitive_scheme
--- PASS: TestExtractBearerJWT/case-insensitive_scheme (0.00s)
=== RUN   TestExtractBearerJWT/virtual_key_is_not_a_jwt
--- PASS: TestExtractBearerJWT/virtual_key_is_not_a_jwt (0.00s)
=== RUN   TestExtractBearerJWT/non-bearer_scheme
--- PASS: TestExtractBearerJWT/non-bearer_scheme (0.00s)
=== RUN   TestExtractBearerJWT/empty_header
--- PASS: TestExtractBearerJWT/empty_header (0.00s)
=== RUN   TestExtractBearerJWT/bearer_without_token
--- PASS: TestExtractBearerJWT/bearer_without_token (0.00s)
--- PASS: TestExtractBearerJWT (0.00s)
=== RUN   TestVerifyMCPJWT_ValidEachMode
=== RUN   TestVerifyMCPJWT_ValidEachMode/user
--- PASS: TestVerifyMCPJWT_ValidEachMode/user (0.00s)
=== RUN   TestVerifyMCPJWT_ValidEachMode/vk
--- PASS: TestVerifyMCPJWT_ValidEachMode/vk (0.00s)
=== RUN   TestVerifyMCPJWT_ValidEachMode/session
--- PASS: TestVerifyMCPJWT_ValidEachMode/session (0.00s)
--- PASS: TestVerifyMCPJWT_ValidEachMode (0.09s)
=== RUN   TestVerifyMCPJWT_Rejections
--- PASS: TestPoolExpiredConnection (1.50s)
=== RUN   TestPoolGetSkipsStaleIdleConnection
--- PASS: TestPoolGetSkipsStaleIdleConnection (0.00s)
=== RUN   TestPoolGetDialsThroughConfiguredHTTPProxy
--- PASS: TestPoolGetDialsThroughConfiguredHTTPProxy (0.00s)
=== RUN   TestPoolGetFailsWithUnreachableProxy
--- PASS: TestPoolGetFailsWithUnreachableProxy (0.00s)
=== RUN   TestSessionManagerCreateAndGet
--- PASS: TestSessionManagerCreateAndGet (0.00s)
=== RUN   TestSessionManagerConnectionLimit
--- PASS: TestSessionManagerConnectionLimit (0.00s)
=== RUN   TestSessionManagerRemove
--- PASS: TestSessionManagerRemove (0.00s)
=== RUN   TestSessionLastResponseID
--- PASS: TestSessionLastResponseID (0.00s)
=== RUN   TestSessionManagerCloseAll
--- PASS: TestSessionManagerCloseAll (0.00s)
=== RUN   TestSessionRealtimeState
--- PASS: TestSessionRealtimeState (0.00s)
=== RUN   TestSessionRecordRealtimeInputUpdatesPendingItemAndIgnoresConsumedLateUpdate
--- PASS: TestSessionRecordRealtimeInputUpdatesPendingItemAndIgnoresConsumedLateUpdate (0.00s)
PASS
ok  	github.com/maximhq/bifrost/transports/bifrost-http/websocket	1.752s
=== RUN   TestVerifyMCPJWT_Rejections/expired_token
--- PASS: TestVerifyMCPJWT_Rejections/expired_token (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/nbf_in_the_future
--- PASS: TestVerifyMCPJWT_Rejections/nbf_in_the_future (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/missing_exp
--- PASS: TestVerifyMCPJWT_Rejections/missing_exp (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/missing_iat
--- PASS: TestVerifyMCPJWT_Rejections/missing_iat (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/issuer_mismatch
--- PASS: TestVerifyMCPJWT_Rejections/issuer_mismatch (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/audience_mismatch
--- PASS: TestVerifyMCPJWT_Rejections/audience_mismatch (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/unknown_kid
--- PASS: TestVerifyMCPJWT_Rejections/unknown_kid (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/wrong_signing_key
--- PASS: TestVerifyMCPJWT_Rejections/wrong_signing_key (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/non-RS256_algorithm_(HS256)
--- PASS: TestVerifyMCPJWT_Rejections/non-RS256_algorithm_(HS256) (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/RS-family_but_not_RS256_(RS384)
--- PASS: TestVerifyMCPJWT_Rejections/RS-family_but_not_RS256_(RS384) (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/alg_none
--- PASS: TestVerifyMCPJWT_Rejections/alg_none (0.00s)
=== RUN   TestVerifyMCPJWT_Rejections/malformed_garbage
--- PASS: TestVerifyMCPJWT_Rejections/malformed_garbage (0.00s)
--- PASS: TestVerifyMCPJWT_Rejections (0.19s)
=== RUN   TestVerifyMCPJWT_NilSigningKeyNotLabeledInvalidToken
--- PASS: TestVerifyMCPJWT_NilSigningKeyNotLabeledInvalidToken (0.16s)
=== RUN   TestCachedSigningKey_ConfigFaults
=== RUN   TestCachedSigningKey_ConfigFaults/nil_config_store
--- PASS: TestCachedSigningKey_ConfigFaults/nil_config_store (0.00s)
=== RUN   TestCachedSigningKey_ConfigFaults/signing_key_load_error
--- PASS: TestCachedSigningKey_ConfigFaults/signing_key_load_error (0.00s)
--- PASS: TestCachedSigningKey_ConfigFaults (0.00s)
=== RUN   TestInjectJWTContext
=== RUN   TestInjectJWTContext/user_mode_sets_user_id
--- PASS: TestInjectJWTContext/user_mode_sets_user_id (0.00s)
=== RUN   TestInjectJWTContext/vk_mode_sets_the_raw_vk_value_and_lets_governance_derive_the_id
--- PASS: TestInjectJWTContext/vk_mode_sets_the_raw_vk_value_and_lets_governance_derive_the_id (0.00s)
=== RUN   TestInjectJWTContext/vk_mode_without_vk_errors
--- PASS: TestInjectJWTContext/vk_mode_without_vk_errors (0.00s)
=== RUN   TestInjectJWTContext/session_mode_sets_session_id
--- PASS: TestInjectJWTContext/session_mode_sets_session_id (0.00s)
=== RUN   TestInjectJWTContext/missing_sub_errors
--- PASS: TestInjectJWTContext/missing_sub_errors (0.00s)
=== RUN   TestInjectJWTContext/unknown_bf_mode_errors
--- PASS: TestInjectJWTContext/unknown_bf_mode_errors (0.00s)
--- PASS: TestInjectJWTContext (0.00s)
=== RUN   TestListSessions
=== RUN   TestListSessions/returns_the_grant_rows
--- PASS: TestListSessions/returns_the_grant_rows (0.00s)
=== RUN   TestListSessions/store_error_surfaces_500
--- PASS: TestListSessions/store_error_surfaces_500 (0.00s)
--- PASS: TestListSessions (0.00s)
=== RUN   TestRevokeSession
=== RUN   TestRevokeSession/vk-mode_grant_revokes_without_an_identity_gate
--- PASS: TestRevokeSession/vk-mode_grant_revokes_without_an_identity_gate (0.00s)
=== RUN   TestRevokeSession/user-mode_grant_revokes_when_the_caller_matches_bf_sub
--- PASS: TestRevokeSession/user-mode_grant_revokes_when_the_caller_matches_bf_sub (0.00s)
=== RUN   TestRevokeSession/revokes_a_visible_grant_regardless_of_caller_identity
--- PASS: TestRevokeSession/revokes_a_visible_grant_regardless_of_caller_identity (0.00s)
=== RUN   TestRevokeSession/not-visible_grant_returns_404_without_attempting_revoke
--- PASS: TestRevokeSession/not-visible_grant_returns_404_without_attempting_revoke (0.00s)
=== RUN   TestRevokeSession/empty_id_returns_400
--- PASS: TestRevokeSession/empty_id_returns_400 (0.00s)
--- PASS: TestRevokeSession (0.00s)
=== RUN   TestMatchRedirectURI
=== RUN   TestMatchRedirectURI/exact_non-loopback_match
--- PASS: TestMatchRedirectURI/exact_non-loopback_match (0.00s)
=== RUN   TestMatchRedirectURI/non-loopback_mismatch
--- PASS: TestMatchRedirectURI/non-loopback_mismatch (0.00s)
=== RUN   TestMatchRedirectURI/loopback_any_port_(127.0.0.1)
--- PASS: TestMatchRedirectURI/loopback_any_port_(127.0.0.1) (0.00s)
=== RUN   TestMatchRedirectURI/loopback_any_port_(localhost)
--- PASS: TestMatchRedirectURI/loopback_any_port_(localhost) (0.00s)
=== RUN   TestMatchRedirectURI/loopback_path_must_still_match
--- PASS: TestMatchRedirectURI/loopback_path_must_still_match (0.00s)
=== RUN   TestMatchRedirectURI/loopback_scheme_must_still_match
--- PASS: TestMatchRedirectURI/loopback_scheme_must_still_match (0.00s)
=== RUN   TestMatchRedirectURI/malformed_candidate
--- PASS: TestMatchRedirectURI/malformed_candidate (0.00s)
=== RUN   TestMatchRedirectURI/no_registered_uris
--- PASS: TestMatchRedirectURI/no_registered_uris (0.00s)
--- PASS: TestMatchRedirectURI (0.00s)
=== RUN   TestOAuth2IssuerURL
=== RUN   TestOAuth2IssuerURL/uses_the_configured_issuer_when_set
--- PASS: TestOAuth2IssuerURL/uses_the_configured_issuer_when_set (0.00s)
=== RUN   TestOAuth2IssuerURL/falls_back_to_the_request_host_when_issuer_is_unset
--- PASS: TestOAuth2IssuerURL/falls_back_to_the_request_host_when_issuer_is_unset (0.00s)
--- PASS: TestOAuth2IssuerURL (0.00s)
=== RUN   TestOAuth2ServerCfg_DefaultsWhenUnset
--- PASS: TestOAuth2ServerCfg_DefaultsWhenUnset (0.00s)
=== RUN   TestRejectStdioMCPClientIfAuthBypassed_UnauthenticatedRejected
--- PASS: TestRejectStdioMCPClientIfAuthBypassed_UnauthenticatedRejected (0.00s)
=== RUN   TestRejectStdioMCPClientIfAuthBypassed_AuthenticatedAllowed
--- PASS: TestRejectStdioMCPClientIfAuthBypassed_AuthenticatedAllowed (0.00s)
=== RUN   TestRejectStdioMCPClientIfAuthBypassed_HTTPUnaffected
--- PASS: TestRejectStdioMCPClientIfAuthBypassed_HTTPUnaffected (0.00s)
=== RUN   TestRejectPrivateMCPTargetIfAuthBypassed_UnauthenticatedLoopbackRejected
--- PASS: TestRejectPrivateMCPTargetIfAuthBypassed_UnauthenticatedLoopbackRejected (0.00s)
=== RUN   TestRejectPrivateMCPTargetIfAuthBypassed_UnauthenticatedLinkLocalRejected
--- PASS: TestRejectPrivateMCPTargetIfAuthBypassed_UnauthenticatedLinkLocalRejected (0.00s)
=== RUN   TestRejectPrivateMCPTargetIfAuthBypassed_AuthenticatedLoopbackAllowed
--- PASS: TestRejectPrivateMCPTargetIfAuthBypassed_AuthenticatedLoopbackAllowed (0.00s)
=== RUN   TestRejectPrivateMCPTargetIfAuthBypassed_UnauthenticatedPublicTargetAllowed
--- PASS: TestRejectPrivateMCPTargetIfAuthBypassed_UnauthenticatedPublicTargetAllowed (0.00s)
=== RUN   TestRejectPrivateMCPTargetIfAuthBypassed_STDIOUnaffected
--- PASS: TestRejectPrivateMCPTargetIfAuthBypassed_STDIOUnaffected (0.00s)
=== RUN   TestGetVKFromRequest
=== RUN   TestGetVKFromRequest/x-bf-vk_header
--- PASS: TestGetVKFromRequest/x-bf-vk_header (0.00s)
=== RUN   TestGetVKFromRequest/Authorization_Bearer_header
--- PASS: TestGetVKFromRequest/Authorization_Bearer_header (0.00s)
=== RUN   TestGetVKFromRequest/x-api-key_header
--- PASS: TestGetVKFromRequest/x-api-key_header (0.00s)
=== RUN   TestGetVKFromRequest/x-goog-api-key_header
--- PASS: TestGetVKFromRequest/x-goog-api-key_header (0.00s)
=== RUN   TestGetVKFromRequest/no_header_returns_empty_string
--- PASS: TestGetVKFromRequest/no_header_returns_empty_string (0.00s)
=== RUN   TestGetVKFromRequest/non-VK_Bearer_token_returns_empty_string
--- PASS: TestGetVKFromRequest/non-VK_Bearer_token_returns_empty_string (0.00s)
=== RUN   TestGetVKFromRequest/non-VK_x-goog-api-key_returns_empty_string
--- PASS: TestGetVKFromRequest/non-VK_x-goog-api-key_returns_empty_string (0.00s)
=== RUN   TestGetVKFromRequest/x-bf-vk_takes_priority_over_x-goog-api-key
--- PASS: TestGetVKFromRequest/x-bf-vk_takes_priority_over_x-goog-api-key (0.00s)
--- PASS: TestGetVKFromRequest (0.00s)
=== RUN   TestAuthenticate_JWTPath
=== RUN   TestAuthenticate_JWTPath/oauth_mode:_vk_JWT_stamps_the_key's_value_for_governance_to_resolve
--- PASS: TestAuthenticate_JWTPath/oauth_mode:_vk_JWT_stamps_the_key's_value_for_governance_to_resolve (0.00s)
=== RUN   TestAuthenticate_JWTPath/vk_JWT_with_an_inactive_key_is_stamped_rather_than_refused_here
--- PASS: TestAuthenticate_JWTPath/vk_JWT_with_an_inactive_key_is_stamped_rather_than_refused_here (0.00s)
=== RUN   TestAuthenticate_JWTPath/vk_JWT_for_an_unknown_key_is_rejected
--- PASS: TestAuthenticate_JWTPath/vk_JWT_for_an_unknown_key_is_rejected (0.00s)
=== RUN   TestAuthenticate_JWTPath/vk_JWT_is_rejected_once_virtual-key_identity_is_disabled_and_user_mode_is_offered
--- PASS: TestAuthenticate_JWTPath/vk_JWT_is_rejected_once_virtual-key_identity_is_disabled_and_user_mode_is_offered (0.00s)
=== RUN   TestAuthenticate_JWTPath/user_JWT_without_a_session_stamps_the_user
--- PASS: TestAuthenticate_JWTPath/user_JWT_without_a_session_stamps_the_user (0.00s)
=== RUN   TestAuthenticate_JWTPath/user_JWT_is_rejected_when_the_user_is_no_longer_active
--- PASS: TestAuthenticate_JWTPath/user_JWT_is_rejected_when_the_user_is_no_longer_active (0.00s)
=== RUN   TestAuthenticate_JWTPath/user_JWT_with_a_matching_session_is_accepted
--- PASS: TestAuthenticate_JWTPath/user_JWT_with_a_matching_session_is_accepted (0.00s)
=== RUN   TestAuthenticate_JWTPath/user_JWT_with_a_mismatched_session_is_rejected
--- PASS: TestAuthenticate_JWTPath/user_JWT_with_a_mismatched_session_is_rejected (0.00s)
=== RUN   TestAuthenticate_JWTPath/both_mode:_a_bearer_an_upstream_auth_layer_authenticated_is_accepted_as_the_stamped_user
--- PASS: TestAuthenticate_JWTPath/both_mode:_a_bearer_an_upstream_auth_layer_authenticated_is_accepted_as_the_stamped_user (0.00s)
=== RUN   TestAuthenticate_JWTPath/both_mode:_a_bearer_no_upstream_auth_layer_authenticated_is_refused_on_its_key_id
--- PASS: TestAuthenticate_JWTPath/both_mode:_a_bearer_no_upstream_auth_layer_authenticated_is_refused_on_its_key_id (0.00s)
=== RUN   TestAuthenticate_JWTPath/oauth_strict_mode_verifies_the_bearer_even_when_an_identity_is_stamped_upstream
--- PASS: TestAuthenticate_JWTPath/oauth_strict_mode_verifies_the_bearer_even_when_an_identity_is_stamped_upstream (0.00s)
=== RUN   TestAuthenticate_JWTPath/session_JWT_is_rejected_when_auth_is_enforced
--- PASS: TestAuthenticate_JWTPath/session_JWT_is_rejected_when_auth_is_enforced (0.00s)
=== RUN   TestAuthenticate_JWTPath/session_JWT_stamps_the_session_when_auth_is_not_enforced
--- PASS: TestAuthenticate_JWTPath/session_JWT_stamps_the_session_when_auth_is_not_enforced (0.00s)
=== RUN   TestAuthenticate_JWTPath/both_mode:_session_token_with_a_header_VK_is_rejected
--- PASS: TestAuthenticate_JWTPath/both_mode:_session_token_with_a_header_VK_is_rejected (0.00s)
=== RUN   TestAuthenticate_JWTPath/both_mode:_vk_token_with_a_header_VK_is_rejected
--- PASS: TestAuthenticate_JWTPath/both_mode:_vk_token_with_a_header_VK_is_rejected (0.00s)
--- PASS: TestAuthenticate_JWTPath (0.08s)
=== RUN   TestAuthenticate_HeaderAndAnonPath
=== RUN   TestAuthenticate_HeaderAndAnonPath/headers_mode:_a_header_VK_is_accepted_without_being_looked_up
--- PASS: TestAuthenticate_HeaderAndAnonPath/headers_mode:_a_header_VK_is_accepted_without_being_looked_up (0.00s)
=== RUN   TestAuthenticate_HeaderAndAnonPath/anonymous_access_is_accepted_when_auth_is_not_enforced
--- PASS: TestAuthenticate_HeaderAndAnonPath/anonymous_access_is_accepted_when_auth_is_not_enforced (0.00s)
=== RUN   TestAuthenticate_HeaderAndAnonPath/no_credentials_rejected_when_auth_is_enforced
--- PASS: TestAuthenticate_HeaderAndAnonPath/no_credentials_rejected_when_auth_is_enforced (0.00s)
=== RUN   TestAuthenticate_HeaderAndAnonPath/headers_mode:_an_identity_stamped_upstream_satisfies_enforced_auth
--- PASS: TestAuthenticate_HeaderAndAnonPath/headers_mode:_an_identity_stamped_upstream_satisfies_enforced_auth (0.00s)
=== RUN   TestAuthenticate_HeaderAndAnonPath/both_mode:_an_identity_stamped_upstream_satisfies_enforced_auth
--- PASS: TestAuthenticate_HeaderAndAnonPath/both_mode:_an_identity_stamped_upstream_satisfies_enforced_auth (0.00s)
=== RUN   TestAuthenticate_HeaderAndAnonPath/oauth_strict_mode_ignores_an_identity_stamped_upstream
--- PASS: TestAuthenticate_HeaderAndAnonPath/oauth_strict_mode_ignores_an_identity_stamped_upstream (0.00s)
=== RUN   TestAuthenticate_HeaderAndAnonPath/oauth_strict_mode_rejects_a_header_VK_with_WWW-Authenticate
--- PASS: TestAuthenticate_HeaderAndAnonPath/oauth_strict_mode_rejects_a_header_VK_with_WWW-Authenticate (0.00s)
=== RUN   TestAuthenticate_HeaderAndAnonPath/oauth_strict_mode_with_no_credentials_sets_WWW-Authenticate
--- PASS: TestAuthenticate_HeaderAndAnonPath/oauth_strict_mode_with_no_credentials_sets_WWW-Authenticate (0.00s)
=== RUN   TestAuthenticate_HeaderAndAnonPath/headers_mode:_a_JWT_bearer_is_not_treated_as_a_credential
--- PASS: TestAuthenticate_HeaderAndAnonPath/headers_mode:_a_JWT_bearer_is_not_treated_as_a_credential (0.00s)
--- PASS: TestAuthenticate_HeaderAndAnonPath (0.13s)
=== RUN   TestAdmit
=== RUN   TestAdmit/a_refusal_is_answered_with_the_status_it_carries
--- PASS: TestAdmit/a_refusal_is_answered_with_the_status_it_carries (0.00s)
=== RUN   TestAdmit/a_401_refusal_points_at_the_authorization_server_when_discovery_is_enabled
--- PASS: TestAdmit/a_401_refusal_points_at_the_authorization_server_when_discovery_is_enabled (0.00s)
=== RUN   TestAdmit/a_refusal_without_a_status_is_forbidden
--- PASS: TestAdmit/a_refusal_without_a_status_is_forbidden (0.00s)
=== RUN   TestAdmit/restricted_access_stamps_the_tools_it_grants
--- PASS: TestAdmit/restricted_access_stamps_the_tools_it_grants (0.00s)
=== RUN   TestAdmit/access_granting_no_tool_stamps_an_empty_list,_which_permits_none
--- PASS: TestAdmit/access_granting_no_tool_stamps_an_empty_list,_which_permits_none (0.00s)
=== RUN   TestAdmit/a_caller's_list_narrows_within_the_access_and_cannot_widen_it
--- PASS: TestAdmit/a_caller's_list_narrows_within_the_access_and_cannot_widen_it (0.00s)
=== RUN   TestAdmit/no_access_leaves_the_caller's_list_alone_and_stamps_nothing_of_its_own
--- PASS: TestAdmit/no_access_leaves_the_caller's_list_alone_and_stamps_nothing_of_its_own (0.00s)
=== RUN   TestAdmit/nothing_resolved_stamps_nothing
--- PASS: TestAdmit/nothing_resolved_stamps_nothing (0.00s)
=== RUN   TestAdmit/governance_is_asked_about_the_request's_own_context,_after_authentication
--- PASS: TestAdmit/governance_is_asked_about_the_request's_own_context,_after_authentication (0.00s)
=== RUN   TestAdmit/a_request_that_fails_authentication_is_refused_before_governance_is_asked
--- PASS: TestAdmit/a_request_that_fails_authentication_is_refused_before_governance_is_asked (0.00s)
--- PASS: TestAdmit (0.18s)
=== RUN   TestConvertToolFunctionParametersToMCPInputSchemaPreservesDefs
--- PASS: TestConvertToolFunctionParametersToMCPInputSchemaPreservesDefs (0.00s)
=== RUN   TestConvertToolFunctionParametersToMCPInputSchemaPreservesLegacyDefinitionsAsDefs
--- PASS: TestConvertToolFunctionParametersToMCPInputSchemaPreservesLegacyDefinitionsAsDefs (0.00s)
=== RUN   TestParseMCPSessionsListQuery_Defaults
--- PASS: TestParseMCPSessionsListQuery_Defaults (0.00s)
=== RUN   TestParseMCPSessionsListQuery_AllParams
--- PASS: TestParseMCPSessionsListQuery_AllParams (0.00s)
=== RUN   TestParseMCPSessionsListQuery_LimitCappedAtMax
--- PASS: TestParseMCPSessionsListQuery_LimitCappedAtMax (0.00s)
=== RUN   TestParseMCPSessionsListQuery_InvalidLimit
--- PASS: TestParseMCPSessionsListQuery_InvalidLimit (0.00s)
=== RUN   TestParseMCPSessionsListQuery_InvalidOffset
--- PASS: TestParseMCPSessionsListQuery_InvalidOffset (0.00s)
=== RUN   TestMCPSessionsListQuery_KindAllowed
--- PASS: TestMCPSessionsListQuery_KindAllowed (0.00s)
=== RUN   TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders
=== RUN   TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//anthropic_passthrough/v1/messages/false
--- PASS: TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//anthropic_passthrough/v1/messages/false (0.00s)
=== RUN   TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//anthropic_passthrough/v1/messages/true
--- PASS: TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//anthropic_passthrough/v1/messages/true (0.00s)
=== RUN   TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//chatgpt_passthrough/backend-api/codex/responses/false
--- PASS: TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//chatgpt_passthrough/backend-api/codex/responses/false (0.00s)
=== RUN   TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//chatgpt_passthrough/backend-api/codex/responses/true
--- PASS: TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//chatgpt_passthrough/backend-api/codex/responses/true (0.00s)
=== RUN   TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//openai_passthrough/v1/responses/false
--- PASS: TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//openai_passthrough/v1/responses/false (0.00s)
=== RUN   TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//openai_passthrough/v1/responses/true
--- PASS: TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders//openai_passthrough/v1/responses/true (0.00s)
--- PASS: TestPassthroughMiddlewaresPreserveOnlyOfficialHeaders (0.00s)
=== RUN   TestCorsMiddleware_LocalhostOrigins
=== RUN   TestCorsMiddleware_LocalhostOrigins/http://localhost:3000
--- PASS: TestCorsMiddleware_LocalhostOrigins/http://localhost:3000 (0.00s)
=== RUN   TestCorsMiddleware_LocalhostOrigins/https://localhost:3000
--- PASS: TestCorsMiddleware_LocalhostOrigins/https://localhost:3000 (0.00s)
=== RUN   TestCorsMiddleware_LocalhostOrigins/http://127.0.0.1:8080
--- PASS: TestCorsMiddleware_LocalhostOrigins/http://127.0.0.1:8080 (0.00s)
=== RUN   TestCorsMiddleware_LocalhostOrigins/http://0.0.0.0:5000
--- PASS: TestCorsMiddleware_LocalhostOrigins/http://0.0.0.0:5000 (0.00s)
=== RUN   TestCorsMiddleware_LocalhostOrigins/https://127.0.0.1:3000
--- PASS: TestCorsMiddleware_LocalhostOrigins/https://127.0.0.1:3000 (0.00s)
--- PASS: TestCorsMiddleware_LocalhostOrigins (0.00s)
=== RUN   TestCorsMiddleware_ConfiguredOrigins
--- PASS: TestCorsMiddleware_ConfiguredOrigins (0.00s)
=== RUN   TestCorsMiddleware_NonAllowedOrigins
--- PASS: TestCorsMiddleware_NonAllowedOrigins (0.00s)
=== RUN   TestCorsMiddleware_PreflightAllowedOrigin
--- PASS: TestCorsMiddleware_PreflightAllowedOrigin (0.00s)
=== RUN   TestCorsMiddleware_PreflightNonAllowedOrigin
--- PASS: TestCorsMiddleware_PreflightNonAllowedOrigin (0.00s)
=== RUN   TestCorsMiddleware_PreflightLocalhost
--- PASS: TestCorsMiddleware_PreflightLocalhost (0.00s)
=== RUN   TestCorsMiddleware_NoOriginHeader
--- PASS: TestCorsMiddleware_NoOriginHeader (0.00s)
=== RUN   TestChainMiddlewares_NoMiddlewares
--- PASS: TestChainMiddlewares_NoMiddlewares (0.00s)
=== RUN   TestChainMiddlewares_SingleMiddleware
--- PASS: TestChainMiddlewares_SingleMiddleware (0.00s)
=== RUN   TestChainMiddlewares_MultipleMiddlewares
--- PASS: TestChainMiddlewares_MultipleMiddlewares (0.00s)
=== RUN   TestChainMiddlewares_MiddlewareCanModifyContext
--- PASS: TestChainMiddlewares_MiddlewareCanModifyContext (0.00s)
=== RUN   TestIsInferenceWSEndpoint
--- PASS: TestIsInferenceWSEndpoint (0.00s)
=== RUN   TestIsRealtimeTransportEndpoint
--- PASS: TestIsRealtimeTransportEndpoint (0.00s)
=== RUN   TestChainMiddlewares_ShortCircuit
--- PASS: TestChainMiddlewares_ShortCircuit (0.00s)
=== RUN   TestChainMiddlewares_ShortCircuitMiddlePosition
--- PASS: TestChainMiddlewares_ShortCircuitMiddlePosition (0.00s)
=== RUN   TestAuthMiddleware_NilAuthConfig
--- PASS: TestAuthMiddleware_NilAuthConfig (0.00s)
=== RUN   TestAuthMiddleware_DisabledAuthConfig
--- PASS: TestAuthMiddleware_DisabledAuthConfig (0.00s)
=== RUN   TestAuthMiddleware_EnabledAuthConfig_NoAuth
--- PASS: TestAuthMiddleware_EnabledAuthConfig_NoAuth (0.00s)
=== RUN   TestAuthMiddleware_SkillsPublicServeManagementSplit
=== RUN   TestAuthMiddleware_SkillsPublicServeManagementSplit/serve_routes_bypass_auth
--- PASS: TestAuthMiddleware_SkillsPublicServeManagementSplit/serve_routes_bypass_auth (0.00s)
=== RUN   TestAuthMiddleware_SkillsPublicServeManagementSplit/management_routes_require_auth
--- PASS: TestAuthMiddleware_SkillsPublicServeManagementSplit/management_routes_require_auth (0.00s)
--- PASS: TestAuthMiddleware_SkillsPublicServeManagementSplit (0.00s)
=== RUN   TestAuthMiddleware_WhitelistedRoutes
=== RUN   TestAuthMiddleware_WhitelistedRoutes//api/session/is-auth-enabled
--- PASS: TestAuthMiddleware_WhitelistedRoutes//api/session/is-auth-enabled (0.00s)
=== RUN   TestAuthMiddleware_WhitelistedRoutes//api/session/login
--- PASS: TestAuthMiddleware_WhitelistedRoutes//api/session/login (0.00s)
=== RUN   TestAuthMiddleware_WhitelistedRoutes//api/session/logout
--- PASS: TestAuthMiddleware_WhitelistedRoutes//api/session/logout (0.00s)
=== RUN   TestAuthMiddleware_WhitelistedRoutes//api/oauth/callback
--- PASS: TestAuthMiddleware_WhitelistedRoutes//api/oauth/callback (0.00s)
=== RUN   TestAuthMiddleware_WhitelistedRoutes//health
--- PASS: TestAuthMiddleware_WhitelistedRoutes//health (0.00s)
--- PASS: TestAuthMiddleware_WhitelistedRoutes (0.00s)
=== RUN   TestAuthMiddleware_APIMiddleware_DevPrefixDoesNotMatchDevices
=== RUN   TestAuthMiddleware_APIMiddleware_DevPrefixDoesNotMatchDevices/dev_pprof_is_whitelisted
--- PASS: TestAuthMiddleware_APIMiddleware_DevPrefixDoesNotMatchDevices/dev_pprof_is_whitelisted (0.00s)
=== RUN   TestAuthMiddleware_APIMiddleware_DevPrefixDoesNotMatchDevices/dev_pprof_subpath_is_whitelisted
--- PASS: TestAuthMiddleware_APIMiddleware_DevPrefixDoesNotMatchDevices/dev_pprof_subpath_is_whitelisted (0.00s)
=== RUN   TestAuthMiddleware_APIMiddleware_DevPrefixDoesNotMatchDevices/devices_is_NOT_whitelisted
--- PASS: TestAuthMiddleware_APIMiddleware_DevPrefixDoesNotMatchDevices/devices_is_NOT_whitelisted (0.00s)
--- PASS: TestAuthMiddleware_APIMiddleware_DevPrefixDoesNotMatchDevices (0.00s)
=== RUN   TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth
=== RUN   TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth//v1/realtime
--- PASS: TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth//v1/realtime (0.00s)
=== RUN   TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth//openai/v1/realtime
--- PASS: TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth//openai/v1/realtime (0.00s)
=== RUN   TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth//v1/realtime/calls?model=gpt-realtime
--- PASS: TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth//v1/realtime/calls?model=gpt-realtime (0.00s)
=== RUN   TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth//openai/v1/realtime/calls?model=gpt-realtime
--- PASS: TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth//openai/v1/realtime/calls?model=gpt-realtime (0.00s)
--- PASS: TestAuthMiddleware_InferenceMiddleware_RealtimeTransportBypassesAuth (0.00s)
=== RUN   TestAuthMiddleware_InferenceMiddleware_DelegatesAuthToGovernance
=== RUN   TestAuthMiddleware_InferenceMiddleware_DelegatesAuthToGovernance/chat_completion_with_virtual_key
--- PASS: TestAuthMiddleware_InferenceMiddleware_DelegatesAuthToGovernance/chat_completion_with_virtual_key (0.00s)
=== RUN   TestAuthMiddleware_InferenceMiddleware_DelegatesAuthToGovernance/chat_completion_without_credentials
--- PASS: TestAuthMiddleware_InferenceMiddleware_DelegatesAuthToGovernance/chat_completion_without_credentials (0.00s)
=== RUN   TestAuthMiddleware_InferenceMiddleware_DelegatesAuthToGovernance/realtime_minting_(client_secrets)
--- PASS: TestAuthMiddleware_InferenceMiddleware_DelegatesAuthToGovernance/realtime_minting_(client_secrets) (0.00s)
--- PASS: TestAuthMiddleware_InferenceMiddleware_DelegatesAuthToGovernance (0.00s)
=== RUN   TestAuthMiddleware_APIMiddleware_VirtualKeyDoesNotBypass
--- PASS: TestAuthMiddleware_APIMiddleware_VirtualKeyDoesNotBypass (0.00s)
=== RUN   TestAuthMiddleware_UpdateAuthConfig_NilToEnabled
--- PASS: TestAuthMiddleware_UpdateAuthConfig_NilToEnabled (0.00s)
=== RUN   TestAuthMiddleware_BootstrapToken_NoAdminNoToken
--- PASS: TestAuthMiddleware_BootstrapToken_NoAdminNoToken (0.00s)
=== RUN   TestAuthMiddleware_BootstrapToken_AdminExists
--- PASS: TestAuthMiddleware_BootstrapToken_AdminExists (0.00s)
=== RUN   TestAuthMiddleware_BootstrapToken_ValidatesAndClears
--- PASS: TestAuthMiddleware_BootstrapToken_ValidatesAndClears (0.00s)
=== RUN   TestAuthMiddleware_UpdateAuthConfig_EnabledToDisabled
--- PASS: TestAuthMiddleware_UpdateAuthConfig_EnabledToDisabled (0.00s)
=== RUN   TestFasthttpToHTTPRequest
--- PASS: TestFasthttpToHTTPRequest (0.00s)
=== RUN   TestCorsMiddleware_DefaultHeaders
--- PASS: TestCorsMiddleware_DefaultHeaders (0.00s)
=== RUN   TestCorsMiddleware_WildcardHeaders_NonCredentialed
--- PASS: TestCorsMiddleware_WildcardHeaders_NonCredentialed (0.00s)
=== RUN   TestCorsMiddleware_WildcardHeaders_CredentialedPreflight
--- PASS: TestCorsMiddleware_WildcardHeaders_CredentialedPreflight (0.00s)
=== RUN   TestCorsMiddleware_WildcardHeaders_CredentialedNonPreflight
--- PASS: TestCorsMiddleware_WildcardHeaders_CredentialedNonPreflight (0.00s)
=== RUN   TestCorsMiddleware_CustomHeaders
--- PASS: TestCorsMiddleware_CustomHeaders (0.00s)
=== RUN   TestCorsMiddleware_DuplicateHeaders
--- PASS: TestCorsMiddleware_DuplicateHeaders (0.00s)
=== RUN   TestCorsMiddleware_CustomHeadersWithLocalhost
--- PASS: TestCorsMiddleware_CustomHeadersWithLocalhost (0.00s)
=== RUN   TestCorsMiddleware_CustomHeadersNotSetForNonAllowedOrigin
--- PASS: TestCorsMiddleware_CustomHeadersNotSetForNonAllowedOrigin (0.00s)
=== RUN   TestFasthttpToHTTPRequest_PathParams
--- PASS: TestFasthttpToHTTPRequest_PathParams (0.00s)
=== RUN   TestRequestDecompressionMiddleware_SupportedEncodings
=== RUN   TestRequestDecompressionMiddleware_SupportedEncodings/gzip
--- PASS: TestRequestDecompressionMiddleware_SupportedEncodings/gzip (0.00s)
=== RUN   TestRequestDecompressionMiddleware_SupportedEncodings/deflate
--- PASS: TestRequestDecompressionMiddleware_SupportedEncodings/deflate (0.00s)
=== RUN   TestRequestDecompressionMiddleware_SupportedEncodings/brotli
--- PASS: TestRequestDecompressionMiddleware_SupportedEncodings/brotli (0.00s)
=== RUN   TestRequestDecompressionMiddleware_SupportedEncodings/zstd
--- PASS: TestRequestDecompressionMiddleware_SupportedEncodings/zstd (0.00s)
--- PASS: TestRequestDecompressionMiddleware_SupportedEncodings (0.00s)
=== RUN   TestRequestDecompressionMiddleware_InvalidCompressedBody
--- PASS: TestRequestDecompressionMiddleware_InvalidCompressedBody (0.00s)
=== RUN   TestRequestDecompressionMiddleware_UnsupportedEncoding
--- PASS: TestRequestDecompressionMiddleware_UnsupportedEncoding (0.00s)
=== RUN   TestRequestDecompressionMiddleware_DecompressedSizeLimit
--- PASS: TestRequestDecompressionMiddleware_DecompressedSizeLimit (0.00s)
=== RUN   TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding
=== RUN   TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding/gzip
--- PASS: TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding/gzip (0.00s)
=== RUN   TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding/deflate
--- PASS: TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding/deflate (0.00s)
=== RUN   TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding/br
--- PASS: TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding/br (0.00s)
=== RUN   TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding/zstd
--- PASS: TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding/zstd (0.00s)
--- PASS: TestRequestDecompressionMiddleware_EmptyBodyWithContentEncoding (0.00s)
=== RUN   TestRequestDecompressionMiddleware_NoContentEncoding
--- PASS: TestRequestDecompressionMiddleware_NoContentEncoding (0.00s)
=== RUN   TestRequestDecompressionMiddleware_ExactSizeLimit
--- PASS: TestRequestDecompressionMiddleware_ExactSizeLimit (0.00s)
=== RUN   TestShouldStreamDecompress
=== RUN   TestShouldStreamDecompress/chunked_(CL=-1)
--- PASS: TestShouldStreamDecompress/chunked_(CL=-1) (0.00s)
=== RUN   TestShouldStreamDecompress/empty_body_(CL=0)
--- PASS: TestShouldStreamDecompress/empty_body_(CL=0) (0.00s)
=== RUN   TestShouldStreamDecompress/small_body
--- PASS: TestShouldStreamDecompress/small_body (0.00s)
=== RUN   TestShouldStreamDecompress/at_default_threshold
--- PASS: TestShouldStreamDecompress/at_default_threshold (0.00s)
=== RUN   TestShouldStreamDecompress/above_default_threshold
--- PASS: TestShouldStreamDecompress/above_default_threshold (0.00s)
=== RUN   TestShouldStreamDecompress/above_custom_threshold
--- PASS: TestShouldStreamDecompress/above_custom_threshold (0.00s)
=== RUN   TestShouldStreamDecompress/below_custom_threshold
--- PASS: TestShouldStreamDecompress/below_custom_threshold (0.00s)
=== RUN   TestShouldStreamDecompress/chunked_with_custom_threshold
--- PASS: TestShouldStreamDecompress/chunked_with_custom_threshold (0.00s)
--- PASS: TestShouldStreamDecompress (0.00s)
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_ChunkedGzip
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_ChunkedGzip (0.00s)
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_AllEncodings
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_AllEncodings/gzip
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_AllEncodings/gzip (0.00s)
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_AllEncodings/deflate
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_AllEncodings/deflate (0.00s)
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_AllEncodings/brotli
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_AllEncodings/brotli (0.00s)
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_AllEncodings/zstd
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_AllEncodings/zstd (0.00s)
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_AllEncodings (0.00s)
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_InvalidBody
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_InvalidBody (0.00s)
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_UnsupportedEncoding
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_UnsupportedEncoding (0.00s)
=== RUN   TestRequestDecompressionMiddleware_BufferedPath_SmallGzip
--- PASS: TestRequestDecompressionMiddleware_BufferedPath_SmallGzip (0.00s)
=== RUN   TestRequestDecompressionMiddleware_StreamingPath_LargeGzip
--- PASS: TestRequestDecompressionMiddleware_StreamingPath_LargeGzip (0.04s)
=== RUN   TestTracingMiddleware_StreamingRootSpanEndsAfterLLMSpan
--- PASS: TestTracingMiddleware_StreamingRootSpanEndsAfterLLMSpan (0.01s)
=== RUN   TestCollectDimensionHeaders
--- PASS: TestCollectDimensionHeaders (0.00s)
=== RUN   TestTracingMiddleware_SetsCorrelationHeaders
=== RUN   TestTracingMiddleware_SetsCorrelationHeaders/generates_request_id_when_absent
--- PASS: TestTracingMiddleware_SetsCorrelationHeaders/generates_request_id_when_absent (0.00s)
=== RUN   TestTracingMiddleware_SetsCorrelationHeaders/echoes_caller-supplied_request_id
--- PASS: TestTracingMiddleware_SetsCorrelationHeaders/echoes_caller-supplied_request_id (0.00s)
=== RUN   TestTracingMiddleware_SetsCorrelationHeaders/headers_survive_the_error_path
--- PASS: TestTracingMiddleware_SetsCorrelationHeaders/headers_survive_the_error_path (0.00s)
--- PASS: TestTracingMiddleware_SetsCorrelationHeaders (0.00s)
=== RUN   TestTracingMiddleware_AccessLogIncludesRequestID
--- PASS: TestTracingMiddleware_AccessLogIncludesRequestID (0.00s)
=== RUN   TestTransportPreAuthInterceptorMiddleware_HeaderVisibleToNext
--- PASS: TestTransportPreAuthInterceptorMiddleware_HeaderVisibleToNext (0.00s)
=== RUN   TestTransportPreAuthInterceptorMiddleware_BodyIsModifiable
--- PASS: TestTransportPreAuthInterceptorMiddleware_BodyIsModifiable (0.00s)
=== RUN   TestTransportPreAuthInterceptorMiddleware_ShortCircuitResponse
--- PASS: TestTransportPreAuthInterceptorMiddleware_ShortCircuitResponse (0.00s)
=== RUN   TestTransportPreAuthInterceptorMiddleware_PluginError
--- PASS: TestTransportPreAuthInterceptorMiddleware_PluginError (0.00s)
=== RUN   TestTransportPreAuthInterceptorMiddleware_PathMutationRejected
--- PASS: TestTransportPreAuthInterceptorMiddleware_PathMutationRejected (0.00s)
=== RUN   TestTransportPreAuthInterceptorMiddleware_NoPlugins
--- PASS: TestTransportPreAuthInterceptorMiddleware_NoPlugins (0.00s)
=== RUN   TestSecurityHeadersMiddleware_APINoStore
=== RUN   TestSecurityHeadersMiddleware_APINoStore/api_path_gets_no-store
--- PASS: TestSecurityHeadersMiddleware_APINoStore/api_path_gets_no-store (0.00s)
=== RUN   TestSecurityHeadersMiddleware_APINoStore/api_path_keeps_handler_policy
--- PASS: TestSecurityHeadersMiddleware_APINoStore/api_path_keeps_handler_policy (0.00s)
=== RUN   TestSecurityHeadersMiddleware_APINoStore/non-api_path_untouched
--- PASS: TestSecurityHeadersMiddleware_APINoStore/non-api_path_untouched (0.00s)
=== RUN   TestSecurityHeadersMiddleware_APINoStore/non-api_path_keeps_handler_policy
--- PASS: TestSecurityHeadersMiddleware_APINoStore/non-api_path_keeps_handler_policy (0.00s)
=== RUN   TestSecurityHeadersMiddleware_APINoStore/prefix_must_match_a_segment
--- PASS: TestSecurityHeadersMiddleware_APINoStore/prefix_must_match_a_segment (0.00s)
--- PASS: TestSecurityHeadersMiddleware_APINoStore (0.00s)
=== RUN   TestGetModelParameters_ResolvesQualifiedAndBareIDs
    model_parameters_test.go:28: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/model_parameters_test.go:28
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestGetModelParameters_ResolvesQualifiedAndBareIDs
--- FAIL: TestGetModelParameters_ResolvesQualifiedAndBareIDs (0.00s)
=== RUN   TestNotificationListFiltersByRole
--- PASS: TestNotificationListFiltersByRole (0.01s)
=== RUN   TestNotificationServicePublishValidatesAndPersists
--- PASS: TestNotificationServicePublishValidatesAndPersists (0.00s)
=== RUN   TestNotificationListStopsAtScanBudget
--- PASS: TestNotificationListStopsAtScanBudget (0.00s)
=== RUN   TestNotificationCreateRequiresLocalAdmin
--- PASS: TestNotificationCreateRequiresLocalAdmin (0.00s)
=== RUN   TestNotificationVisibleToRole
--- PASS: TestNotificationVisibleToRole (0.00s)
=== RUN   TestNotificationCursorRoundTrip
--- PASS: TestNotificationCursorRoundTrip (0.00s)
=== RUN   TestCreatePlugin_RejectsCustomPathWhenAuthBypassed
--- PASS: TestCreatePlugin_RejectsCustomPathWhenAuthBypassed (0.00s)
=== RUN   TestCreatePlugin_AllowsCustomPathWhenNotBypassed
--- PASS: TestCreatePlugin_AllowsCustomPathWhenNotBypassed (0.00s)
=== RUN   TestUpdatePlugin_RejectsCustomPathWhenAuthBypassed
--- PASS: TestUpdatePlugin_RejectsCustomPathWhenAuthBypassed (0.00s)
=== RUN   TestRestoreRedacted_OTELProfilesHeaders
--- PASS: TestRestoreRedacted_OTELProfilesHeaders (0.00s)
=== RUN   TestRestoreRedacted_KafkaSecretVarObjects
--- PASS: TestRestoreRedacted_KafkaSecretVarObjects (0.00s)
=== RUN   TestRestoreRedacted_FullyRedactedSentinel
--- PASS: TestRestoreRedacted_FullyRedactedSentinel (0.00s)
=== RUN   TestUpdatePlugin_ConfigMerge
--- PASS: TestUpdatePlugin_ConfigMerge (0.00s)
=== RUN   TestUpdatePlugin_ConfigMerge_NewPlugin
--- PASS: TestUpdatePlugin_ConfigMerge_NewPlugin (0.00s)
=== RUN   TestGetLoadedPlugins
--- PASS: TestGetLoadedPlugins (0.00s)
=== RUN   TestUpdatePricingOverride_ReplacesFullBody
    pricing_override_test.go:105: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:105
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestUpdatePricingOverride_ReplacesFullBody
--- FAIL: TestUpdatePricingOverride_ReplacesFullBody (0.00s)
=== RUN   TestPromptMutationsReloadCache
=== RUN   TestPromptMutationsReloadCache/deleteFolder
--- PASS: TestPromptMutationsReloadCache/deleteFolder (0.00s)
=== RUN   TestPromptMutationsReloadCache/createPrompt
--- PASS: TestPromptMutationsReloadCache/createPrompt (0.00s)
=== RUN   TestPromptMutationsReloadCache/updatePrompt
--- PASS: TestPromptMutationsReloadCache/updatePrompt (0.00s)
=== RUN   TestPromptMutationsReloadCache/deletePrompt
--- PASS: TestPromptMutationsReloadCache/deletePrompt (0.00s)
=== RUN   TestPromptMutationsReloadCache/createVersion
--- PASS: TestPromptMutationsReloadCache/createVersion (0.00s)
=== RUN   TestPromptMutationsReloadCache/deleteVersion
--- PASS: TestPromptMutationsReloadCache/deleteVersion (0.00s)
=== RUN   TestPromptMutationsReloadCache/createSession
--- PASS: TestPromptMutationsReloadCache/createSession (0.00s)
=== RUN   TestPromptMutationsReloadCache/updateSession
--- PASS: TestPromptMutationsReloadCache/updateSession (0.00s)
=== RUN   TestPromptMutationsReloadCache/deleteSession
--- PASS: TestPromptMutationsReloadCache/deleteSession (0.00s)
=== RUN   TestPromptMutationsReloadCache/renameSession
--- PASS: TestPromptMutationsReloadCache/renameSession (0.00s)
=== RUN   TestPromptMutationsReloadCache/commitSession
--- PASS: TestPromptMutationsReloadCache/commitSession (0.00s)
--- PASS: TestPromptMutationsReloadCache (0.01s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/deleteFolder
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/deleteFolder (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/createPrompt
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/createPrompt (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/updatePrompt
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/updatePrompt (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/deletePrompt
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/deletePrompt (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/createVersion
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/createVersion (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/deleteVersion
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/deleteVersion (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/createSession
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/createSession (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/updateSession
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/updateSession (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/deleteSession
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/deleteSession (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/renameSession
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/renameSession (0.00s)
=== RUN   TestPromptMutationsDoNotReloadOnStoreFailure/commitSession
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure/commitSession (0.00s)
--- PASS: TestPromptMutationsDoNotReloadOnStoreFailure (0.00s)
=== RUN   TestPromptMutationsSurviveReloaderProblems
=== RUN   TestPromptMutationsSurviveReloaderProblems/nil_reloader
--- PASS: TestPromptMutationsSurviveReloaderProblems/nil_reloader (0.00s)
=== RUN   TestPromptMutationsSurviveReloaderProblems/failing_reloader
--- PASS: TestPromptMutationsSurviveReloaderProblems/failing_reloader (0.00s)
--- PASS: TestPromptMutationsSurviveReloaderProblems (0.00s)
=== RUN   TestDatabricksKey_RedactedUpdateRoundTrip
    provider_keys_databricks_test.go:36: redacted workspace_url="https://dbc-1.cloud.databricks.com" api_format="ai_gateway" tags=true client_id="dbx-************************alue"
--- PASS: TestDatabricksKey_RedactedUpdateRoundTrip (0.00s)
=== RUN   TestDatabricksKey_Validate
=== RUN   TestDatabricksKey_Validate/pat_ok
--- PASS: TestDatabricksKey_Validate/pat_ok (0.00s)
=== RUN   TestDatabricksKey_Validate/oauth_ok
--- PASS: TestDatabricksKey_Validate/oauth_ok (0.00s)
=== RUN   TestDatabricksKey_Validate/nil_config
--- PASS: TestDatabricksKey_Validate/nil_config (0.00s)
=== RUN   TestDatabricksKey_Validate/no_workspace_url
--- PASS: TestDatabricksKey_Validate/no_workspace_url (0.00s)
=== RUN   TestDatabricksKey_Validate/no_auth_at_all
--- PASS: TestDatabricksKey_Validate/no_auth_at_all (0.00s)
=== RUN   TestDatabricksKey_Validate/half_a_service_principal
--- PASS: TestDatabricksKey_Validate/half_a_service_principal (0.00s)
=== RUN   TestDatabricksKey_Validate/half_a_pair_alongside_a_pat
--- PASS: TestDatabricksKey_Validate/half_a_pair_alongside_a_pat (0.00s)
=== RUN   TestDatabricksKey_Validate/bad_api_format
--- PASS: TestDatabricksKey_Validate/bad_api_format (0.00s)
=== RUN   TestDatabricksKey_Validate/env-ref_workspace_url
--- PASS: TestDatabricksKey_Validate/env-ref_workspace_url (0.00s)
=== RUN   TestDatabricksKey_Validate/env-ref_pat
--- PASS: TestDatabricksKey_Validate/env-ref_pat (0.00s)
--- PASS: TestDatabricksKey_Validate (0.00s)
=== RUN   TestDatabricksKey_SwitchToPAT
--- PASS: TestDatabricksKey_SwitchToPAT (0.00s)
=== RUN   TestMergeUpdatedKey_Value
=== RUN   TestMergeUpdatedKey_Value/echoed_current_redaction_preserves_stored_value
--- PASS: TestMergeUpdatedKey_Value/echoed_current_redaction_preserves_stored_value (0.00s)
=== RUN   TestMergeUpdatedKey_Value/mismatched_mask_still_preserves_stored_value
--- PASS: TestMergeUpdatedKey_Value/mismatched_mask_still_preserves_stored_value (0.00s)
=== RUN   TestMergeUpdatedKey_Value/genuine_new_plaintext_value_is_applied
--- PASS: TestMergeUpdatedKey_Value/genuine_new_plaintext_value_is_applied (0.00s)
=== RUN   TestMergeUpdatedKey_Value/genuine_env_ref_is_applied_not_preserved
--- PASS: TestMergeUpdatedKey_Value/genuine_env_ref_is_applied_not_preserved (0.00s)
=== RUN   TestMergeUpdatedKey_Value/empty_value_is_not_treated_as_redacted
--- PASS: TestMergeUpdatedKey_Value/empty_value_is_not_treated_as_redacted (0.00s)
--- PASS: TestMergeUpdatedKey_Value (0.00s)
=== RUN   TestMergeUpdatedKey_Name
=== RUN   TestMergeUpdatedKey_Name/name_omitted_from_update_preserves_stored_name
--- PASS: TestMergeUpdatedKey_Name/name_omitted_from_update_preserves_stored_name (0.00s)
=== RUN   TestMergeUpdatedKey_Name/explicit_new_name_is_applied
--- PASS: TestMergeUpdatedKey_Name/explicit_new_name_is_applied (0.00s)
--- PASS: TestMergeUpdatedKey_Name (0.00s)
=== RUN   TestMergeUpdatedKey_ProviderConfigMaskedPreviews
--- PASS: TestMergeUpdatedKey_ProviderConfigMaskedPreviews (0.00s)
=== RUN   TestMergeUpdatedKey_GithubCopilotKeyConfig
=== RUN   TestMergeUpdatedKey_GithubCopilotKeyConfig/masks_resolve_back_to_the_stored_values
--- PASS: TestMergeUpdatedKey_GithubCopilotKeyConfig/masks_resolve_back_to_the_stored_values (0.11s)
=== RUN   TestMergeUpdatedKey_GithubCopilotKeyConfig/a_mask_with_no_stored_counterpart_is_rejected
--- PASS: TestMergeUpdatedKey_GithubCopilotKeyConfig/a_mask_with_no_stored_counterpart_is_rejected (0.00s)
=== RUN   TestMergeUpdatedKey_GithubCopilotKeyConfig/unmasked_literals_overwrite_the_stored_values
--- PASS: TestMergeUpdatedKey_GithubCopilotKeyConfig/unmasked_literals_overwrite_the_stored_values (0.00s)
--- PASS: TestMergeUpdatedKey_GithubCopilotKeyConfig (0.11s)
=== RUN   TestMergeUpdatedKey_RejectsMaskWithoutStoredCounterpart
=== RUN   TestMergeUpdatedKey_RejectsMaskWithoutStoredCounterpart/missing_config_section
--- PASS: TestMergeUpdatedKey_RejectsMaskWithoutStoredCounterpart/missing_config_section (0.00s)
=== RUN   TestMergeUpdatedKey_RejectsMaskWithoutStoredCounterpart/missing_optional_field
--- PASS: TestMergeUpdatedKey_RejectsMaskWithoutStoredCounterpart/missing_optional_field (0.00s)
=== RUN   TestMergeUpdatedKey_RejectsMaskWithoutStoredCounterpart/empty_stored_value
--- PASS: TestMergeUpdatedKey_RejectsMaskWithoutStoredCounterpart/empty_stored_value (0.00s)
--- PASS: TestMergeUpdatedKey_RejectsMaskWithoutStoredCounterpart (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields
=== RUN   TestValidateProviderKeyRequiredNestedFields/azure_missing_endpoint
--- PASS: TestValidateProviderKeyRequiredNestedFields/azure_missing_endpoint (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/azure_nil_config
--- PASS: TestValidateProviderKeyRequiredNestedFields/azure_nil_config (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/azure_ok
--- PASS: TestValidateProviderKeyRequiredNestedFields/azure_ok (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/bedrock_missing_region
--- PASS: TestValidateProviderKeyRequiredNestedFields/bedrock_missing_region (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/bedrock_ok
--- PASS: TestValidateProviderKeyRequiredNestedFields/bedrock_ok (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/mantle_missing_region
--- PASS: TestValidateProviderKeyRequiredNestedFields/mantle_missing_region (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/mantle_ok
--- PASS: TestValidateProviderKeyRequiredNestedFields/mantle_ok (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/vllm_missing_url
--- PASS: TestValidateProviderKeyRequiredNestedFields/vllm_missing_url (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/vllm_missing_model_name
--- PASS: TestValidateProviderKeyRequiredNestedFields/vllm_missing_model_name (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/vllm_ok
--- PASS: TestValidateProviderKeyRequiredNestedFields/vllm_ok (0.00s)
=== RUN   TestValidateProviderKeyRequiredNestedFields/openai_unaffected
--- PASS: TestValidateProviderKeyRequiredNestedFields/openai_unaffected (0.00s)
--- PASS: TestValidateProviderKeyRequiredNestedFields (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats
=== RUN   TestValidateProviderKeyGithubCopilotFormats/valid_literal_app_config
--- PASS: TestValidateProviderKeyGithubCopilotFormats/valid_literal_app_config (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/direct_token_needs_no_app_config
--- PASS: TestValidateProviderKeyGithubCopilotFormats/direct_token_needs_no_app_config (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/non-numeric_installation_id
--- PASS: TestValidateProviderKeyGithubCopilotFormats/non-numeric_installation_id (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/path_traversal_in_installation_id
--- PASS: TestValidateProviderKeyGithubCopilotFormats/path_traversal_in_installation_id (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/non-numeric_repository_id
--- PASS: TestValidateProviderKeyGithubCopilotFormats/non-numeric_repository_id (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/private_key_that_is_not_PEM
--- PASS: TestValidateProviderKeyGithubCopilotFormats/private_key_that_is_not_PEM (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/valid_PKCS#8_private_key
--- PASS: TestValidateProviderKeyGithubCopilotFormats/valid_PKCS#8_private_key (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/EC_private_key
--- PASS: TestValidateProviderKeyGithubCopilotFormats/EC_private_key (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/encrypted_private_key_envelope
--- PASS: TestValidateProviderKeyGithubCopilotFormats/encrypted_private_key_envelope (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/well-formed_envelope_with_a_garbage_DER_payload
--- PASS: TestValidateProviderKeyGithubCopilotFormats/well-formed_envelope_with_a_garbage_DER_payload (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/PEM_whose_newlines_survived_as_literal_backslash-n
--- PASS: TestValidateProviderKeyGithubCopilotFormats/PEM_whose_newlines_survived_as_literal_backslash-n (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/no_token_and_no_app_config_at_all
--- PASS: TestValidateProviderKeyGithubCopilotFormats/no_token_and_no_app_config_at_all (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/no_token,_missing_app_id
--- PASS: TestValidateProviderKeyGithubCopilotFormats/no_token,_missing_app_id (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/no_token,_missing_installation_id
--- PASS: TestValidateProviderKeyGithubCopilotFormats/no_token,_missing_installation_id (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/no_token,_missing_repository_id
--- PASS: TestValidateProviderKeyGithubCopilotFormats/no_token,_missing_repository_id (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/no_token,_missing_private_key
--- PASS: TestValidateProviderKeyGithubCopilotFormats/no_token,_missing_private_key (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/direct_token_alongside_an_incomplete_app_config
--- PASS: TestValidateProviderKeyGithubCopilotFormats/direct_token_alongside_an_incomplete_app_config (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/direct_token_alongside_a_malformed_app_config
--- PASS: TestValidateProviderKeyGithubCopilotFormats/direct_token_alongside_a_malformed_app_config (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/direct_token_alongside_a_complete_app_config
--- PASS: TestValidateProviderKeyGithubCopilotFormats/direct_token_alongside_a_complete_app_config (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/client_ID_as_app_id_is_accepted
--- PASS: TestValidateProviderKeyGithubCopilotFormats/client_ID_as_app_id_is_accepted (0.00s)
=== RUN   TestValidateProviderKeyGithubCopilotFormats/env_references_are_not_format-checked
--- PASS: TestValidateProviderKeyGithubCopilotFormats/env_references_are_not_format-checked (0.00s)
--- PASS: TestValidateProviderKeyGithubCopilotFormats (0.00s)
=== RUN   TestRefreshProviderModels_DelegatesToModelsManager
--- PASS: TestRefreshProviderModels_DelegatesToModelsManager (0.01s)
=== RUN   TestRefreshProviderKeyModels_DelegatesToModelsManager
--- PASS: TestRefreshProviderKeyModels_DelegatesToModelsManager (0.00s)
=== RUN   TestRefreshProviderModels_InFlightReturns409
--- PASS: TestRefreshProviderModels_InFlightReturns409 (0.00s)
=== RUN   TestRefreshProviderKeyModels_UnknownKeyReturns404
--- PASS: TestRefreshProviderKeyModels_UnknownKeyReturns404 (0.00s)
=== RUN   TestRefreshProviderModels_UnknownProviderReturns404
--- PASS: TestRefreshProviderModels_UnknownProviderReturns404 (0.00s)
=== RUN   TestCreateProviderKey_CustomBedrockRequiresRegion
--- PASS: TestCreateProviderKey_CustomBedrockRequiresRegion (0.02s)
=== RUN   TestApplyProviderConfigUpdates_OmittedBlocksArePreserved
--- PASS: TestApplyProviderConfigUpdates_OmittedBlocksArePreserved (0.01s)
=== RUN   TestApplyProviderConfigUpdates_ExplicitNullClears
--- PASS: TestApplyProviderConfigUpdates_ExplicitNullClears (0.00s)
=== RUN   TestApplyProviderConfigUpdates_PresentBlocksAreReplaced
--- PASS: TestApplyProviderConfigUpdates_PresentBlocksAreReplaced (0.00s)
=== RUN   TestAddProvider_ReloadsRuntimeEvenWhenModelDiscoveryIsSkipped
--- PASS: TestAddProvider_ReloadsRuntimeEvenWhenModelDiscoveryIsSkipped (0.01s)
=== RUN   TestAddProvider_ReturnsErrorWhenRuntimeReloadFails
--- PASS: TestAddProvider_ReturnsErrorWhenRuntimeReloadFails (0.00s)
=== RUN   TestUpdateProvider_RejectsKeysInBody
--- PASS: TestUpdateProvider_RejectsKeysInBody (0.02s)
=== RUN   TestUpdateProvider_PassesThroughForEmptyOrAbsentKeys
=== RUN   TestUpdateProvider_PassesThroughForEmptyOrAbsentKeys/keys_field_omitted_entirely
--- PASS: TestUpdateProvider_PassesThroughForEmptyOrAbsentKeys/keys_field_omitted_entirely (0.00s)
=== RUN   TestUpdateProvider_PassesThroughForEmptyOrAbsentKeys/keys_explicitly_null
--- PASS: TestUpdateProvider_PassesThroughForEmptyOrAbsentKeys/keys_explicitly_null (0.00s)
=== RUN   TestUpdateProvider_PassesThroughForEmptyOrAbsentKeys/keys_explicitly_empty_array
--- PASS: TestUpdateProvider_PassesThroughForEmptyOrAbsentKeys/keys_explicitly_empty_array (0.00s)
--- PASS: TestUpdateProvider_PassesThroughForEmptyOrAbsentKeys (0.00s)
=== RUN   TestListModels_UnknownKeysDoNotFilter
--- PASS: TestListModels_UnknownKeysDoNotFilter (0.00s)
=== RUN   TestListModels_ReturnsExactAccessibleByKeysAndSkipsDisabledKeys
--- PASS: TestListModels_ReturnsExactAccessibleByKeysAndSkipsDisabledKeys (0.00s)
=== RUN   TestListModels_AppliesQueryAndLimitAfterFiltering
--- PASS: TestListModels_AppliesQueryAndLimitAfterFiltering (0.00s)
=== RUN   TestListModels_MarksDeprecatedModelsWithoutFiltering
    bedrock_test.go:95: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/integrations/bedrock_test.go:95
        	Error:      	Received unexpected error:
        	            	failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestGenericRouter_MarkDeprecatedListModelsResponseUsesCatalog2402190276\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestGenericRouter_MarkDeprecatedListModelsResponseUsesCatalog2402190276\\001\\pricing.json" after host
        	Test:       	TestGenericRouter_MarkDeprecatedListModelsResponseUsesCatalog
--- FAIL: TestGenericRouter_MarkDeprecatedListModelsResponseUsesCatalog (7.00s)
=== RUN   Test_parseS3URI
=== RUN   Test_parseS3URI/full_S3_URI_with_key
--- PASS: Test_parseS3URI/full_S3_URI_with_key (0.00s)
=== RUN   Test_parseS3URI/S3_URI_with_bucket_only
--- PASS: Test_parseS3URI/S3_URI_with_bucket_only (0.00s)
=== RUN   Test_parseS3URI/S3_URI_with_bucket_no_trailing_slash
--- PASS: Test_parseS3URI/S3_URI_with_bucket_no_trailing_slash (0.00s)
=== RUN   Test_parseS3URI/plain_bucket_name
--- PASS: Test_parseS3URI/plain_bucket_name (0.00s)
=== RUN   Test_parseS3URI/S3_URI_with_nested_key
--- PASS: Test_parseS3URI/S3_URI_with_nested_key (0.00s)
=== RUN   Test_parseS3URI/empty_string
--- PASS: Test_parseS3URI/empty_string (0.00s)
--- PASS: Test_parseS3URI (0.00s)
=== RUN   Test_createBedrockRouteConfigs
--- PASS: Test_createBedrockRouteConfigs (0.00s)
=== RUN   Test_createBedrockConverseRouteConfig
--- PASS: Test_createBedrockConverseRouteConfig (0.00s)
=== RUN   Test_createBedrockConverseStreamRouteConfig
--- PASS: Test_createBedrockConverseStreamRouteConfig (0.00s)
=== RUN   Test_createBedrockInvokeRouteConfig
--- PASS: Test_createBedrockInvokeRouteConfig (0.00s)
=== RUN   TestBedrockInvokeEmbeddingResponseLangChainCohereAlias
--- PASS: TestBedrockInvokeEmbeddingResponseLangChainCohereAlias (0.00s)
=== RUN   TestBedrockInvokeTypedCohereResponseLangChainAlias
--- PASS: TestBedrockInvokeTypedCohereResponseLangChainAlias (0.00s)
=== RUN   TestBedrockInvokeTypedCohereResponseWithoutFloatIsNotAliased
--- PASS: TestBedrockInvokeTypedCohereResponseWithoutFloatIsNotAliased (0.00s)
=== RUN   TestBedrockInvokeEmbeddingResponseLangChainLeavesTitanUnchanged
--- PASS: TestBedrockInvokeEmbeddingResponseLangChainLeavesTitanUnchanged (0.00s)
=== RUN   TestBedrockInvokeTypedTitanResponseLangChainLeavesEnvelopeUnchanged
--- PASS: TestBedrockInvokeTypedTitanResponseLangChainLeavesEnvelopeUnchanged (0.00s)
=== RUN   TestBedrockInvokeTypedCohereResponseUsesResolvedModelForAlias
--- PASS: TestBedrockInvokeTypedCohereResponseUsesResolvedModelForAlias (0.00s)
=== RUN   Test_createBedrockInvokeWithResponseStreamRouteConfig
--- PASS: Test_createBedrockInvokeWithResponseStreamRouteConfig (0.00s)
=== RUN   Test_bedrockStreamErrorConverterEncodesEventStreamException
--- PASS: Test_bedrockStreamErrorConverterEncodesEventStreamException (0.00s)
=== RUN   Test_toBedrockEventStreamExceptionAcceptsBedrockError
--- PASS: Test_toBedrockEventStreamExceptionAcceptsBedrockError (0.00s)
=== RUN   Test_handleStreamingBedrockUnknownErrorResponseFallsBackToEventStreamException
--- PASS: Test_handleStreamingBedrockUnknownErrorResponseFallsBackToEventStreamException (0.00s)
=== RUN   Test_createBedrockRerankRouteConfig
--- PASS: Test_createBedrockRerankRouteConfig (0.00s)
=== RUN   Test_createBedrockRerankResponseConverterEmitsBedrockShape
--- PASS: Test_createBedrockRerankResponseConverterEmitsBedrockShape (0.00s)
=== RUN   Test_createBedrockRerankResponseConverterUsesRawResponse
--- PASS: Test_createBedrockRerankResponseConverterUsesRawResponse (0.00s)
=== RUN   Test_createBedrockRerankResponseConverterConvertsForCrossProvider
--- PASS: Test_createBedrockRerankResponseConverterConvertsForCrossProvider (0.00s)
=== RUN   Test_createBedrockRerankRouteRequestConverter
--- PASS: Test_createBedrockRerankRouteRequestConverter (0.00s)
=== RUN   Test_createBedrockRouteConfigsIncludesRerankForCompositePrefixes
--- PASS: Test_createBedrockRouteConfigsIncludesRerankForCompositePrefixes (0.00s)
=== RUN   Test_createBedrockBatchRouteConfigs
--- PASS: Test_createBedrockBatchRouteConfigs (0.00s)
=== RUN   Test_createBedrockFilesRouteConfigs
--- PASS: Test_createBedrockFilesRouteConfigs (0.00s)
=== RUN   Test_parseS3PutObjectRequest
=== RUN   Test_parseS3PutObjectRequest/valid_request
--- PASS: Test_parseS3PutObjectRequest/valid_request (0.00s)
=== RUN   Test_parseS3PutObjectRequest/simple_key_without_folder
--- PASS: Test_parseS3PutObjectRequest/simple_key_without_folder (0.00s)
=== RUN   Test_parseS3PutObjectRequest/missing_bucket
--- PASS: Test_parseS3PutObjectRequest/missing_bucket (0.00s)
=== RUN   Test_parseS3PutObjectRequest/missing_key
--- PASS: Test_parseS3PutObjectRequest/missing_key (0.00s)
--- PASS: Test_parseS3PutObjectRequest (0.00s)
=== RUN   Test_parseS3PutObjectRequest_invalidType
--- PASS: Test_parseS3PutObjectRequest_invalidType (0.00s)
=== RUN   Test_s3PutObjectPostCallback
=== RUN   Test_s3PutObjectPostCallback/valid_response_with_ID
--- PASS: Test_s3PutObjectPostCallback/valid_response_with_ID (0.00s)
=== RUN   Test_s3PutObjectPostCallback/nil_response
--- PASS: Test_s3PutObjectPostCallback/nil_response (0.00s)
--- PASS: Test_s3PutObjectPostCallback (0.00s)
=== RUN   Test_s3GetObjectPostCallback
=== RUN   Test_s3GetObjectPostCallback/valid_response
--- PASS: Test_s3GetObjectPostCallback/valid_response (0.00s)
=== RUN   Test_s3GetObjectPostCallback/nil_response
--- PASS: Test_s3GetObjectPostCallback/nil_response (0.00s)
--- PASS: Test_s3GetObjectPostCallback (0.00s)
=== RUN   Test_s3HeadObjectPostCallback
=== RUN   Test_s3HeadObjectPostCallback/valid_response
--- PASS: Test_s3HeadObjectPostCallback/valid_response (0.00s)
=== RUN   Test_s3HeadObjectPostCallback/nil_response
--- PASS: Test_s3HeadObjectPostCallback/nil_response (0.00s)
--- PASS: Test_s3HeadObjectPostCallback (0.00s)
=== RUN   Test_s3DeleteObjectPostCallback
--- PASS: Test_s3DeleteObjectPostCallback (0.00s)
=== RUN   Test_s3ListObjectsV2PostCallback
--- PASS: Test_s3ListObjectsV2PostCallback (0.00s)
=== RUN   Test_extractBedrockBatchListQueryParams
=== RUN   Test_extractBedrockBatchListQueryParams/all_params
--- PASS: Test_extractBedrockBatchListQueryParams/all_params (0.00s)
=== RUN   Test_extractBedrockBatchListQueryParams/no_params
--- PASS: Test_extractBedrockBatchListQueryParams/no_params (0.00s)
=== RUN   Test_extractBedrockBatchListQueryParams/invalid_maxResults
--- PASS: Test_extractBedrockBatchListQueryParams/invalid_maxResults (0.00s)
--- PASS: Test_extractBedrockBatchListQueryParams (0.00s)
=== RUN   Test_extractBedrockJobArnFromPath
=== RUN   Test_extractBedrockJobArnFromPath/valid_job_ARN_for_Bedrock
--- PASS: Test_extractBedrockJobArnFromPath/valid_job_ARN_for_Bedrock (0.00s)
=== RUN   Test_extractBedrockJobArnFromPath/URL_encoded_job_ARN
--- PASS: Test_extractBedrockJobArnFromPath/URL_encoded_job_ARN (0.00s)
=== RUN   Test_extractBedrockJobArnFromPath/non-Bedrock_provider_strips_ARN_prefix
--- PASS: Test_extractBedrockJobArnFromPath/non-Bedrock_provider_strips_ARN_prefix (0.00s)
=== RUN   Test_extractBedrockJobArnFromPath/missing_job_arn
--- PASS: Test_extractBedrockJobArnFromPath/missing_job_arn (0.00s)
=== RUN   Test_extractBedrockJobArnFromPath/empty_job_arn
--- PASS: Test_extractBedrockJobArnFromPath/empty_job_arn (0.00s)
--- PASS: Test_extractBedrockJobArnFromPath (0.00s)
=== RUN   Test_extractS3ListObjectsV2Params
=== RUN   Test_extractS3ListObjectsV2Params/all_params
--- PASS: Test_extractS3ListObjectsV2Params/all_params (0.00s)
=== RUN   Test_extractS3ListObjectsV2Params/bucket_only
--- PASS: Test_extractS3ListObjectsV2Params/bucket_only (0.00s)
=== RUN   Test_extractS3ListObjectsV2Params/missing_bucket
--- PASS: Test_extractS3ListObjectsV2Params/missing_bucket (0.00s)
--- PASS: Test_extractS3ListObjectsV2Params (0.00s)
=== RUN   Test_extractS3BucketKeyFromPath
=== RUN   Test_extractS3BucketKeyFromPath/content_operation
--- PASS: Test_extractS3BucketKeyFromPath/content_operation (0.00s)
=== RUN   Test_extractS3BucketKeyFromPath/missing_bucket
--- PASS: Test_extractS3BucketKeyFromPath/missing_bucket (0.00s)
=== RUN   Test_extractS3BucketKeyFromPath/missing_key
--- PASS: Test_extractS3BucketKeyFromPath/missing_key (0.00s)
--- PASS: Test_extractS3BucketKeyFromPath (0.00s)
=== RUN   TestCreateCohereRouteConfigsIncludesRerank
--- PASS: TestCreateCohereRouteConfigsIncludesRerank (0.00s)
=== RUN   TestCohereRerankRouteRequestConverter
--- PASS: TestCohereRerankRouteRequestConverter (0.00s)
=== RUN   TestCohereRerankResponseConverterEmitsCohereShape
--- PASS: TestCohereRerankResponseConverterEmitsCohereShape (0.00s)
=== RUN   TestCohereRerankResponseConverterUsesRawResponse
--- PASS: TestCohereRerankResponseConverterUsesRawResponse (0.00s)
=== RUN   TestCohereRerankResponseConverterConvertsForCrossProvider
--- PASS: TestCohereRerankResponseConverterConvertsForCrossProvider (0.00s)
=== RUN   TestCohereChatResponseConverterEmitsCohereV2Shape
--- PASS: TestCohereChatResponseConverterEmitsCohereV2Shape (0.00s)
=== RUN   TestCohereChatResponseConverterUsesNativeRawResponse
--- PASS: TestCohereChatResponseConverterUsesNativeRawResponse (0.00s)
=== RUN   TestCohereChatResponseConverterConvertsCrossProviderRawResponse
--- PASS: TestCohereChatResponseConverterConvertsCrossProviderRawResponse (0.00s)
=== RUN   TestExtractModelAndRequestType_LargePayloadUsesMetadataWithoutBodyParse
--- PASS: TestExtractModelAndRequestType_LargePayloadUsesMetadataWithoutBodyParse (0.00s)
=== RUN   TestExtractModelAndRequestType_LargeBodyHeuristicSkipsParse
--- PASS: TestExtractModelAndRequestType_LargeBodyHeuristicSkipsParse (0.00s)
=== RUN   TestRewriteGenAIRawRequestBodyRedactsOnlyContentFields
--- PASS: TestRewriteGenAIRawRequestBodyRedactsOnlyContentFields (0.00s)
=== RUN   TestRewriteGenAIRawRequestBodySupportsCountTokensEnvelope
--- PASS: TestRewriteGenAIRawRequestBodySupportsCountTokensEnvelope (0.00s)
=== RUN   TestRewriteGenAIRawRequestBodyRejectsUnmappedLiteral
--- PASS: TestRewriteGenAIRawRequestBodyRejectsUnmappedLiteral (0.00s)
=== RUN   TestRewriteGenAIRawRequestBodyRejectsSensitiveObjectKeys
--- PASS: TestRewriteGenAIRawRequestBodyRejectsSensitiveObjectKeys (0.00s)
=== RUN   TestCreateGenAIRerankRouteConfig
--- PASS: TestCreateGenAIRerankRouteConfig (0.00s)
=== RUN   TestCreateGenAIRouteConfigsIncludesRerank
--- PASS: TestCreateGenAIRouteConfigsIncludesRerank (0.00s)
=== RUN   TestGenAISpeechStreamResponseConverter
--- PASS: TestGenAISpeechStreamResponseConverter (0.00s)
=== RUN   TestGenAISpeechStreamDoneResponseIncludesUsageAndFinishReason
--- PASS: TestGenAISpeechStreamDoneResponseIncludesUsageAndFinishReason (0.00s)
=== RUN   TestGenAICamelCaseSpeechVoiceSurvivesConversion
--- PASS: TestGenAICamelCaseSpeechVoiceSurvivesConversion (0.05s)
=== RUN   TestExtractAndSetModelAndRequestTypePreservesRawBodyForGenerateContent
--- PASS: TestExtractAndSetModelAndRequestTypePreservesRawBodyForGenerateContent (0.00s)
=== RUN   TestExtractAndSetModelAndRequestTypeNoRawPassthroughWithoutExplicitGemini
--- PASS: TestExtractAndSetModelAndRequestTypeNoRawPassthroughWithoutExplicitGemini (0.00s)
=== RUN   TestExtractAndSetModelAndRequestTypeDoesNotRawPassthroughEmbedding
--- PASS: TestExtractAndSetModelAndRequestTypeDoesNotRawPassthroughEmbedding (0.00s)
=== RUN   TestExtractAndSetModelAndRequestTypeDetectsImageEditAtAnyPartPosition
=== RUN   TestExtractAndSetModelAndRequestTypeDetectsImageEditAtAnyPartPosition/text_first
--- PASS: TestExtractAndSetModelAndRequestTypeDetectsImageEditAtAnyPartPosition/text_first (0.00s)
=== RUN   TestExtractAndSetModelAndRequestTypeDetectsImageEditAtAnyPartPosition/image_first
--- PASS: TestExtractAndSetModelAndRequestTypeDetectsImageEditAtAnyPartPosition/image_first (0.00s)
--- PASS: TestExtractAndSetModelAndRequestTypeDetectsImageEditAtAnyPartPosition (0.00s)
=== RUN   TestGenAIBatchCreateConverterCarriesRawBody
--- PASS: TestGenAIBatchCreateConverterCarriesRawBody (0.02s)
=== RUN   TestGenAIBatchCreateConverterCarriesDisplayNameAndInlineTools
--- PASS: TestGenAIBatchCreateConverterCarriesDisplayNameAndInlineTools (0.03s)
=== RUN   TestGenAICachedContentCreateParserRejectsNonStringScalars
--- PASS: TestGenAICachedContentCreateParserRejectsNonStringScalars (0.00s)
=== RUN   TestGenAICachedContentCreateParserCarriesRawBody
--- PASS: TestGenAICachedContentCreateParserCarriesRawBody (0.00s)
=== RUN   TestCreateGenAIRouteConfigsIncludesRerankForCompositePrefixes
--- PASS: TestCreateGenAIRouteConfigsIncludesRerankForCompositePrefixes (0.00s)
=== RUN   TestGenAIRerankRequestConverter
--- PASS: TestGenAIRerankRequestConverter (0.00s)
=== RUN   TestGenAIRerankResponseConverterRestoresCallerRecordIDs
--- PASS: TestGenAIRerankResponseConverterRestoresCallerRecordIDs (0.00s)
=== RUN   TestGenAIRerankRequestConverterRequestsDocuments
--- PASS: TestGenAIRerankRequestConverterRequestsDocuments (0.00s)
=== RUN   TestCreateGenAIRouteConfigsIncludesModelMetadataRoute
--- PASS: TestCreateGenAIRouteConfigsIncludesModelMetadataRoute (0.00s)
=== RUN   TestExtractGeminiModelMetadataParams
--- PASS: TestExtractGeminiModelMetadataParams (0.00s)
=== RUN   TestConvertGeminiModelMetadataResponse
--- PASS: TestConvertGeminiModelMetadataResponse (0.00s)
=== RUN   TestConvertGeminiModelMetadataResponse_MatchesRequestedModelNotFirst
--- PASS: TestConvertGeminiModelMetadataResponse_MatchesRequestedModelNotFirst (0.00s)
=== RUN   TestConvertGeminiModelMetadataResponse_EmptyReturnsMinimalModel
--- PASS: TestConvertGeminiModelMetadataResponse_EmptyReturnsMinimalModel (0.00s)
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting/gemini_to_gemini/generateContent
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting/gemini_to_gemini/generateContent (0.00s)
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting/gemini_to_gemini/streamGenerateContent
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting/gemini_to_gemini/streamGenerateContent (0.00s)
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting/gemini_to_vertex/generateContent
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting/gemini_to_vertex/generateContent (0.00s)
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting/gemini_to_vertex/streamGenerateContent
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting/gemini_to_vertex/streamGenerateContent (0.00s)
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting/vertex_to_gemini/generateContent
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting/vertex_to_gemini/generateContent (0.00s)
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting/vertex_to_gemini/streamGenerateContent
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting/vertex_to_gemini/streamGenerateContent (0.00s)
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting/vertex_to_vertex/generateContent
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting/vertex_to_vertex/generateContent (0.00s)
=== RUN   TestGenAIFlashLiteMinimalThinkingAfterRouting/vertex_to_vertex/streamGenerateContent
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting/vertex_to_vertex/streamGenerateContent (0.00s)
--- PASS: TestGenAIFlashLiteMinimalThinkingAfterRouting (0.00s)
=== RUN   TestOpenAIWireCostResponse
=== RUN   TestOpenAIWireCostResponse/chat_renders_cost_as_float_total_without_mutating_the_shared_response
--- PASS: TestOpenAIWireCostResponse/chat_renders_cost_as_float_total_without_mutating_the_shared_response (0.00s)
=== RUN   TestOpenAIWireCostResponse/responses_renders_cost_as_float_total
--- PASS: TestOpenAIWireCostResponse/responses_renders_cost_as_float_total (0.01s)
=== RUN   TestOpenAIWireCostResponse/image_renders_cost_as_float_total
--- PASS: TestOpenAIWireCostResponse/image_renders_cost_as_float_total (0.00s)
=== RUN   TestOpenAIWireCostResponse/chat_without_cost_passes_through_the_same_pointer
--- PASS: TestOpenAIWireCostResponse/chat_without_cost_passes_through_the_same_pointer (0.00s)
=== RUN   TestOpenAIWireCostResponse/raw_upstream_payload_passes_through_unchanged
--- PASS: TestOpenAIWireCostResponse/raw_upstream_payload_passes_through_unchanged (0.00s)
--- PASS: TestOpenAIWireCostResponse (0.02s)
=== RUN   TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter
=== RUN   TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/both_fields_present
--- PASS: TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/both_fields_present (0.00s)
=== RUN   TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/neither_field_present
--- PASS: TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/neither_field_present (0.00s)
=== RUN   TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/anchor_only
--- PASS: TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/anchor_only (0.00s)
=== RUN   TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/seconds_only
--- PASS: TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/seconds_only (0.00s)
=== RUN   TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/out-of-range_seconds
--- PASS: TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/out-of-range_seconds (0.00s)
=== RUN   TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/non-numeric_seconds
--- PASS: TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter/non-numeric_seconds (0.00s)
--- PASS: TestParseOpenAIFileUploadMultipartRequest_ExpiresAfter (0.00s)
=== RUN   TestParseTranscriptionMultipartRequest_ExtraParamsPassthrough
--- PASS: TestParseTranscriptionMultipartRequest_ExtraParamsPassthrough (0.00s)
=== RUN   TestParseTranscriptionMultipartRequest_TypedFieldsNotShadowedByExtraParams
--- PASS: TestParseTranscriptionMultipartRequest_TypedFieldsNotShadowedByExtraParams (0.00s)
=== RUN   TestParseTranscriptionMultipartRequest_PreservesFilename
--- PASS: TestParseTranscriptionMultipartRequest_PreservesFilename (0.00s)
=== RUN   TestForwardPassthroughRequestHeaderKeepsCookie
--- PASS: TestForwardPassthroughRequestHeaderKeepsCookie (0.00s)
=== RUN   Test_handleStreamingSSESendsHeartbeatDuringIdleGap
--- PASS: Test_handleStreamingSSESendsHeartbeatDuringIdleGap (0.00s)
=== RUN   Test_handleStreamingGenAIEmitsNoHeartbeat
--- PASS: Test_handleStreamingGenAIEmitsNoHeartbeat (0.00s)
=== RUN   Test_passthroughHeartbeatEligible
=== RUN   Test_passthroughHeartbeatEligible/plain_SSE
--- PASS: Test_passthroughHeartbeatEligible/plain_SSE (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/SSE_with_charset_param
--- PASS: Test_passthroughHeartbeatEligible/SSE_with_charset_param (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/SSE_uppercase
--- PASS: Test_passthroughHeartbeatEligible/SSE_uppercase (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/SSE_with_surrounding_space_before_param
--- PASS: Test_passthroughHeartbeatEligible/SSE_with_surrounding_space_before_param (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/raw_JSON_passthrough
--- PASS: Test_passthroughHeartbeatEligible/raw_JSON_passthrough (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/empty_content-type
--- PASS: Test_passthroughHeartbeatEligible/empty_content-type (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/binary_passthrough
--- PASS: Test_passthroughHeartbeatEligible/binary_passthrough (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/prefix-collision_suffix
--- PASS: Test_passthroughHeartbeatEligible/prefix-collision_suffix (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/prefix-collision_longer_type
--- PASS: Test_passthroughHeartbeatEligible/prefix-collision_longer_type (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/Gemini_SSE_passthrough
--- PASS: Test_passthroughHeartbeatEligible/Gemini_SSE_passthrough (0.00s)
=== RUN   Test_passthroughHeartbeatEligible/Vertex_SSE_passthrough
--- PASS: Test_passthroughHeartbeatEligible/Vertex_SSE_passthrough (0.00s)
--- PASS: Test_passthroughHeartbeatEligible (0.00s)
=== RUN   TestCreateHandler_SkipsRequestParserInLargePayloadMode
--- PASS: TestCreateHandler_SkipsRequestParserInLargePayloadMode (0.00s)
=== RUN   TestCreateHandler_UsesRequestParserWhenNotInLargePayloadMode
--- PASS: TestCreateHandler_UsesRequestParserWhenNotInLargePayloadMode (0.00s)
=== RUN   TestResolveLargePayloadMetadata_NilContext
--- PASS: TestResolveLargePayloadMetadata_NilContext (0.00s)
=== RUN   TestResolveLargePayloadMetadata_SyncPath
--- PASS: TestResolveLargePayloadMetadata_SyncPath (0.00s)
=== RUN   TestResolveLargePayloadMetadata_DeferredReady
--- PASS: TestResolveLargePayloadMetadata_DeferredReady (0.00s)
=== RUN   TestResolveLargePayloadMetadata_DeferredNotReady
--- PASS: TestResolveLargePayloadMetadata_DeferredNotReady (0.00s)
=== RUN   TestResolveLargePayloadMetadata_SyncTakesPrecedence
--- PASS: TestResolveLargePayloadMetadata_SyncTakesPrecedence (0.00s)
=== RUN   TestListModelsNarrowsTheFanOut
--- PASS: TestListModelsNarrowsTheFanOut (0.00s)
=== RUN   TestListModelsLeavesFanOutAloneWhenNothingResolved
--- PASS: TestListModelsLeavesFanOutAloneWhenNothingResolved (0.00s)
=== RUN   TestListModelsWithoutAccessResolver
--- PASS: TestListModelsWithoutAccessResolver (0.00s)
=== RUN   TestGrantedProvidersShapeMatchesWhatTheRouterPublishes
--- PASS: TestGrantedProvidersShapeMatchesWhatTheRouterPublishes (0.00s)
=== RUN   Test_handleStreamingInterceptionErrorUsesErrorConverter
--- PASS: Test_handleStreamingInterceptionErrorUsesErrorConverter (0.00s)
=== RUN   Test_handleStreamingInterceptionErrorAnthropicSSEFraming
--- PASS: Test_handleStreamingInterceptionErrorAnthropicSSEFraming (0.00s)
=== RUN   Test_handleStreamingInterceptionErrorIsSanitized
--- PASS: Test_handleStreamingInterceptionErrorIsSanitized (0.00s)
=== RUN   Test_handleStreamingInterceptionErrorBedrockEventStream
--- PASS: Test_handleStreamingInterceptionErrorBedrockEventStream (0.00s)
=== RUN   Test_handleStreamingPlainInterceptionErrorKeepsFlatFormat
--- PASS: Test_handleStreamingPlainInterceptionErrorKeepsFlatFormat (0.00s)
=== RUN   Test_handleStreamingMissingSpeechConverterReturnsError
--- PASS: Test_handleStreamingMissingSpeechConverterReturnsError (0.00s)
=== RUN   TestIsPassthroughRequestOnlyMatchesRegisteredNativeRoutes
--- PASS: TestIsPassthroughRequestOnlyMatchesRegisteredNativeRoutes (0.00s)
=== RUN   TestParsePassthroughBody_MultipartExtractsModelAfterFilePart
--- PASS: TestParsePassthroughBody_MultipartExtractsModelAfterFilePart (0.00s)
=== RUN   TestChatGPTPassthroughRouterRegistersCodexResponsesPost
--- PASS: TestChatGPTPassthroughRouterRegistersCodexResponsesPost (0.00s)
=== RUN   TestChatGPTBackgroundRoutes
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/codex/models?client_version=0.159.3
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/codex/models?client_version=0.159.3 (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/HEAD/backend-api/codex/models
--- PASS: TestChatGPTBackgroundRoutes/HEAD/backend-api/codex/models (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/ps/plugins/list
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/ps/plugins/list (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/POST/backend-api/ps/plugins/install
--- PASS: TestChatGPTBackgroundRoutes/POST/backend-api/ps/plugins/install (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/ps/mcp
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/ps/mcp (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/PUT/backend-api/ps/mcp/config
--- PASS: TestChatGPTBackgroundRoutes/PUT/backend-api/ps/mcp/config (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/wham/usage
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/wham/usage (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/plugins/featured
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/plugins/featured (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/POST/backend-api/codex/analytics-events/events
--- PASS: TestChatGPTBackgroundRoutes/POST/backend-api/codex/analytics-events/events (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/codex/analytics-events/events
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/codex/analytics-events/events (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/POST/backend-api/codex/models
--- PASS: TestChatGPTBackgroundRoutes/POST/backend-api/codex/models (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/DELETE/backend-api/accounts
--- PASS: TestChatGPTBackgroundRoutes/DELETE/backend-api/accounts (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/ps/mcp-lookalike
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/ps/mcp-lookalike (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/GET/backend-api/ps/plugins-other/list
--- PASS: TestChatGPTBackgroundRoutes/GET/backend-api/ps/plugins-other/list (0.00s)
=== RUN   TestChatGPTBackgroundRoutes/TRACE/backend-api/wham/usage
--- PASS: TestChatGPTBackgroundRoutes/TRACE/backend-api/wham/usage (0.00s)
--- PASS: TestChatGPTBackgroundRoutes (0.00s)
=== RUN   TestChatGPTBackgroundWireAndHooks
--- PASS: TestChatGPTBackgroundWireAndHooks (0.01s)
=== RUN   TestChatGPTUpstreamURLDefaultsToChatGPT
--- PASS: TestChatGPTUpstreamURLDefaultsToChatGPT (0.00s)
=== RUN   TestChatGPTUpstreamURLHonoursDiagnosticOverride
--- PASS: TestChatGPTUpstreamURLHonoursDiagnosticOverride (0.00s)
=== RUN   TestRunwarePassthroughRouterRegistersCatchAll
--- PASS: TestRunwarePassthroughRouterRegistersCatchAll (0.00s)
=== RUN   TestRequestWithSettableExtraParams_OpenAIChatRequest
=== RUN   TestRequestWithSettableExtraParams_OpenAIChatRequest/SetExtraParams_populates_both_standalone_and_embedded_ExtraParams
--- PASS: TestRequestWithSettableExtraParams_OpenAIChatRequest/SetExtraParams_populates_both_standalone_and_embedded_ExtraParams (0.00s)
=== RUN   TestRequestWithSettableExtraParams_OpenAIChatRequest/extra_params_propagate_through_ToBifrostChatRequest
--- PASS: TestRequestWithSettableExtraParams_OpenAIChatRequest/extra_params_propagate_through_ToBifrostChatRequest (0.00s)
--- PASS: TestRequestWithSettableExtraParams_OpenAIChatRequest (0.00s)
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIChatRequest_implements_RequestWithSettableExtraParams
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIChatRequest_implements_RequestWithSettableExtraParams (0.00s)
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAITextCompletionRequest_implements_RequestWithSettableExtraParams
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAITextCompletionRequest_implements_RequestWithSettableExtraParams (0.00s)
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIResponsesRequest_implements_RequestWithSettableExtraParams
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIResponsesRequest_implements_RequestWithSettableExtraParams (0.00s)
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIEmbeddingRequest_implements_RequestWithSettableExtraParams
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIEmbeddingRequest_implements_RequestWithSettableExtraParams (0.00s)
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAISpeechRequest_implements_RequestWithSettableExtraParams
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAISpeechRequest_implements_RequestWithSettableExtraParams (0.00s)
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIImageGenerationRequest_implements_RequestWithSettableExtraParams
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIImageGenerationRequest_implements_RequestWithSettableExtraParams (0.00s)
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIImageEditRequest_implements_RequestWithSettableExtraParams
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIImageEditRequest_implements_RequestWithSettableExtraParams (0.00s)
=== RUN   TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIImageVariationRequest_implements_RequestWithSettableExtraParams
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes/OpenAIImageVariationRequest_implements_RequestWithSettableExtraParams (0.00s)
--- PASS: TestRequestWithSettableExtraParams_AllOpenAIRequestTypes (0.00s)
=== RUN   TestExtraParamsRequiresPassthroughHeader
=== RUN   TestExtraParamsRequiresPassthroughHeader/extra_params_NOT_extracted_without_passthrough_header
--- PASS: TestExtraParamsRequiresPassthroughHeader/extra_params_NOT_extracted_without_passthrough_header (0.03s)
=== RUN   TestExtraParamsRequiresPassthroughHeader/extra_params_extracted_with_passthrough_header
--- PASS: TestExtraParamsRequiresPassthroughHeader/extra_params_extracted_with_passthrough_header (0.00s)
--- PASS: TestExtraParamsRequiresPassthroughHeader (0.04s)
=== RUN   TestExtraParamsPassthrough_NestedStructures
--- PASS: TestExtraParamsPassthrough_NestedStructures (0.00s)
=== RUN   TestExtraParamsPassthrough_EndToEnd
--- PASS: TestExtraParamsPassthrough_EndToEnd (0.00s)
=== RUN   TestExtraParamsPassthrough_NoExtraParamsKey
--- PASS: TestExtraParamsPassthrough_NoExtraParamsKey (0.00s)
=== RUN   TestOpenAIChatStructuredOutputRequestParserAndConverter
--- PASS: TestOpenAIChatStructuredOutputRequestParserAndConverter (0.00s)
=== RUN   TestCreateHandler_AnthropicRouteSetsPassthroughFlags
--- PASS: TestCreateHandler_AnthropicRouteSetsPassthroughFlags (0.00s)
=== RUN   TestCreateHandler_AnthropicThreadContinueRefused
--- PASS: TestCreateHandler_AnthropicThreadContinueRefused (0.00s)
=== RUN   TestCreateHandler_AnthropicThreadContinueRefusedStreaming
--- PASS: TestCreateHandler_AnthropicThreadContinueRefusedStreaming (0.00s)
=== RUN   TestCreateHandler_AnthropicThreadCreateProceeds
--- PASS: TestCreateHandler_AnthropicThreadCreateProceeds (0.00s)
=== RUN   TestCreateHandler_CustomParserFailureClosesConnection
--- PASS: TestCreateHandler_CustomParserFailureClosesConnection (0.00s)
=== RUN   TestCreateHandler_DefaultJSONParserFailureClosesConnection
--- PASS: TestCreateHandler_DefaultJSONParserFailureClosesConnection (0.00s)
=== RUN   TestCreateHandler_ParseFailureClosesKeepAliveSocket
=== RUN   TestCreateHandler_ParseFailureClosesKeepAliveSocket/valid_keep-alive_requests_reuse_the_socket
--- PASS: TestCreateHandler_ParseFailureClosesKeepAliveSocket/valid_keep-alive_requests_reuse_the_socket (0.00s)
=== RUN   TestCreateHandler_ParseFailureClosesKeepAliveSocket/malformed_request_closes_the_socket
--- PASS: TestCreateHandler_ParseFailureClosesKeepAliveSocket/malformed_request_closes_the_socket (0.00s)
--- PASS: TestCreateHandler_ParseFailureClosesKeepAliveSocket (0.00s)
=== RUN   TestExtraParamsSetViaInterfaceMutatesOriginalReq
--- PASS: TestExtraParamsSetViaInterfaceMutatesOriginalReq (0.00s)
=== RUN   TestExtractModelFromPath
=== RUN   TestExtractModelFromPath/azure_deployment_chat
--- PASS: TestExtractModelFromPath/azure_deployment_chat (0.00s)
=== RUN   TestExtractModelFromPath/azure_deployment_leading-stripped
--- PASS: TestExtractModelFromPath/azure_deployment_leading-stripped (0.00s)
=== RUN   TestExtractModelFromPath/genai_models_with_action
--- PASS: TestExtractModelFromPath/genai_models_with_action (0.00s)
=== RUN   TestExtractModelFromPath/genai_models_stream_action
--- PASS: TestExtractModelFromPath/genai_models_stream_action (0.00s)
=== RUN   TestExtractModelFromPath/genai_tunedModels
--- PASS: TestExtractModelFromPath/genai_tunedModels (0.00s)
=== RUN   TestExtractModelFromPath/vertex_fully-qualified
--- PASS: TestExtractModelFromPath/vertex_fully-qualified (0.00s)
=== RUN   TestExtractModelFromPath/no_model_segment
--- PASS: TestExtractModelFromPath/no_model_segment (0.00s)
=== RUN   TestExtractModelFromPath/deployments_with_no_trailing_segment
--- PASS: TestExtractModelFromPath/deployments_with_no_trailing_segment (0.00s)
--- PASS: TestExtractModelFromPath (0.00s)
=== RUN   TestExtractPassthroughModel
=== RUN   TestExtractPassthroughModel/azure_deployment_path_overrides_empty_body
--- PASS: TestExtractPassthroughModel/azure_deployment_path_overrides_empty_body (0.00s)
=== RUN   TestExtractPassthroughModel/body_fallback_when_path_has_no_model
--- PASS: TestExtractPassthroughModel/body_fallback_when_path_has_no_model (0.00s)
=== RUN   TestExtractPassthroughModel/path_wins_over_body
--- PASS: TestExtractPassthroughModel/path_wins_over_body (0.00s)
=== RUN   TestExtractPassthroughModel/both_empty
--- PASS: TestExtractPassthroughModel/both_empty (0.00s)
=== RUN   TestExtractPassthroughModel/vertex_resource_body_model
--- PASS: TestExtractPassthroughModel/vertex_resource_body_model (0.00s)
=== RUN   TestExtractPassthroughModel/genai_resource_body_model
--- PASS: TestExtractPassthroughModel/genai_resource_body_model (0.00s)
=== RUN   TestExtractPassthroughModel/slashed_non-resource_model_untouched
--- PASS: TestExtractPassthroughModel/slashed_non-resource_model_untouched (0.00s)
--- PASS: TestExtractPassthroughModel (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_AnthropicOAuthForwarded
--- PASS: TestApplyPassthroughCallerAuth_AnthropicOAuthForwarded (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_OpenAIJWTForwarded
--- PASS: TestApplyPassthroughCallerAuth_OpenAIJWTForwarded (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_LoopbackHTTPOverrideForwardsJWT
--- PASS: TestApplyPassthroughCallerAuth_LoopbackHTTPOverrideForwardsJWT (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_APIKeysStayStripped
=== RUN   TestApplyPassthroughCallerAuth_APIKeysStayStripped/openai_plain_api_key
--- PASS: TestApplyPassthroughCallerAuth_APIKeysStayStripped/openai_plain_api_key (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_APIKeysStayStripped/anthropic_api_key_bearer
--- PASS: TestApplyPassthroughCallerAuth_APIKeysStayStripped/anthropic_api_key_bearer (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_APIKeysStayStripped/provider_override_bedrock
--- PASS: TestApplyPassthroughCallerAuth_APIKeysStayStripped/provider_override_bedrock (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_APIKeysStayStripped/openai_two-segment_token
--- PASS: TestApplyPassthroughCallerAuth_APIKeysStayStripped/openai_two-segment_token (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_APIKeysStayStripped/http_upstream_override
--- PASS: TestApplyPassthroughCallerAuth_APIKeysStayStripped/http_upstream_override (0.00s)
=== RUN   TestApplyPassthroughCallerAuth_APIKeysStayStripped/http_lookalike_host
--- PASS: TestApplyPassthroughCallerAuth_APIKeysStayStripped/http_lookalike_host (0.00s)
--- PASS: TestApplyPassthroughCallerAuth_APIKeysStayStripped (0.00s)
=== RUN   TestExtractAndParseFallbacks_GeminiGenerationRequest
--- PASS: TestExtractAndParseFallbacks_GeminiGenerationRequest (0.00s)
=== RUN   TestJSONErrorsPreserveBodyForBodylessStatus
--- PASS: TestJSONErrorsPreserveBodyForBodylessStatus (0.00s)
=== RUN   TestSendStreamError_PropagatesProviderStatusCode
=== RUN   TestSendStreamError_PropagatesProviderStatusCode/provider_400_-_Bedrock_ValidationException_/_OpenAI_invalid_request_error
--- PASS: TestSendStreamError_PropagatesProviderStatusCode/provider_400_-_Bedrock_ValidationException_/_OpenAI_invalid_request_error (0.00s)
=== RUN   TestSendStreamError_PropagatesProviderStatusCode/provider_429_-_rate_limiting_(all_providers)
--- PASS: TestSendStreamError_PropagatesProviderStatusCode/provider_429_-_rate_limiting_(all_providers) (0.00s)
=== RUN   TestSendStreamError_PropagatesProviderStatusCode/provider_503_-_Bedrock_ServiceUnavailableException
--- PASS: TestSendStreamError_PropagatesProviderStatusCode/provider_503_-_Bedrock_ServiceUnavailableException (0.00s)
=== RUN   TestSendStreamError_PropagatesProviderStatusCode/provider_529_-_Anthropic_overloaded_error
--- PASS: TestSendStreamError_PropagatesProviderStatusCode/provider_529_-_Anthropic_overloaded_error (0.00s)
=== RUN   TestSendStreamError_PropagatesProviderStatusCode/nil_StatusCode_defaults_to_500
--- PASS: TestSendStreamError_PropagatesProviderStatusCode/nil_StatusCode_defaults_to_500 (0.00s)
--- PASS: TestSendStreamError_PropagatesProviderStatusCode (0.00s)
=== RUN   TestSendStreamError_OpenAIErrorFormat
--- PASS: TestSendStreamError_OpenAIErrorFormat (0.00s)
=== RUN   TestSendStreamError_AnthropicErrorFormat
--- PASS: TestSendStreamError_AnthropicErrorFormat (0.00s)
=== RUN   TestSendStreamError_BedrockErrorFormat
--- PASS: TestSendStreamError_BedrockErrorFormat (0.00s)
=== RUN   TestSendStreamError_ForwardsProviderHeaders
--- PASS: TestSendStreamError_ForwardsProviderHeaders (0.00s)
=== RUN   TestTryStreamLargeResponse_AppliesRoutedIdentityHeaders
--- PASS: TestTryStreamLargeResponse_AppliesRoutedIdentityHeaders (0.00s)
FAIL
FAIL	github.com/maximhq/bifrost/transports/bifrost-http/integrations	7.473s
    providers_test.go:549: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModels_MarksDeprecatedModelsWithoutFiltering2590695793\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModels_MarksDeprecatedModelsWithoutFiltering2590695793\\001\\pricing.json" after host
--- FAIL: TestListModels_MarksDeprecatedModelsWithoutFiltering (7.00s)
=== RUN   TestListBaseModels_IncludesDeprecatedPricingRows
    providers_test.go:592: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListBaseModels_IncludesDeprecatedPricingRows1786824924\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListBaseModels_IncludesDeprecatedPricingRows1786824924\\001\\pricing.json" after host
--- FAIL: TestListBaseModels_IncludesDeprecatedPricingRows (7.00s)
=== RUN   TestEnrichListModelsResponse_MarksDeprecatedPricingRows
    providers_test.go:617: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestEnrichListModelsResponse_MarksDeprecatedPricingRows1431118774\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestEnrichListModelsResponse_MarksDeprecatedPricingRows1431118774\\001\\pricing.json" after host
--- FAIL: TestEnrichListModelsResponse_MarksDeprecatedPricingRows (7.00s)
=== RUN   TestListModels_UnfilteredIgnoresKeys
--- PASS: TestListModels_UnfilteredIgnoresKeys (0.00s)
=== RUN   TestListModels_UnfilteredWithoutKeysReturnsAllUnfilteredModels
--- PASS: TestListModels_UnfilteredWithoutKeysReturnsAllUnfilteredModels (0.00s)
=== RUN   TestListModelDetails_ErrorsWhenModelCatalogUnavailable
--- PASS: TestListModelDetails_ErrorsWhenModelCatalogUnavailable (0.00s)
=== RUN   TestListModelDetails_UnknownKeysDoNotFilter
--- PASS: TestListModelDetails_UnknownKeysDoNotFilter (0.00s)
=== RUN   TestListModelDetails_SkipsUnknownKeysAndFiltersWithValid
--- PASS: TestListModelDetails_SkipsUnknownKeysAndFiltersWithValid (0.00s)
=== RUN   TestListModelDetails_SkipsDisabledKeysAndFiltersWithValid
--- PASS: TestListModelDetails_SkipsDisabledKeysAndFiltersWithValid (0.00s)
=== RUN   TestListModelDetails_UnfilteredIgnoresKeys
--- PASS: TestListModelDetails_UnfilteredIgnoresKeys (0.00s)
=== RUN   TestListModelDetails_IncludesPricing
    providers_test.go:888: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_IncludesPricing3497550844\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_IncludesPricing3497550844\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_IncludesPricing (7.00s)
=== RUN   TestListModelDetails_ResolvesCatalogPricing
=== RUN   TestListModelDetails_ResolvesCatalogPricing/Together_catalog_provider
    providers_test.go:1027: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingTogether_catalog_pro1675286047\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingTogether_catalog_pro1675286047\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_ResolvesCatalogPricing/Together_catalog_provider (7.00s)
=== RUN   TestListModelDetails_ResolvesCatalogPricing/Together_alias
    providers_test.go:1027: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingTogether_alias1648311547\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingTogether_alias1648311547\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_ResolvesCatalogPricing/Together_alias (7.00s)
=== RUN   TestListModelDetails_ResolvesCatalogPricing/Azure_alias
    providers_test.go:1027: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingAzure_alias81113808\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingAzure_alias81113808\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_ResolvesCatalogPricing/Azure_alias (7.00s)
=== RUN   TestListModelDetails_ResolvesCatalogPricing/Azure_alias_with_empty_model_name
    providers_test.go:1027: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingAzure_alias_with_emp136421870\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingAzure_alias_with_emp136421870\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_ResolvesCatalogPricing/Azure_alias_with_empty_model_name (7.00s)
=== RUN   TestListModelDetails_ResolvesCatalogPricing/Azure_alias_falls_back_to_alias_key
    providers_test.go:1027: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingAzure_alias_falls_ba3209452414\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_ResolvesCatalogPricingAzure_alias_falls_ba3209452414\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_ResolvesCatalogPricing/Azure_alias_falls_back_to_alias_key (7.00s)
--- FAIL: TestListModelDetails_ResolvesCatalogPricing (35.02s)
=== RUN   TestListModelDetails_AppliesGlobalOverrideWithoutMutatingBase
    providers_test.go:1094: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_AppliesGlobalOverrideWithoutMutatingBase3592389\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_AppliesGlobalOverrideWithoutMutatingBase3592389\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_AppliesGlobalOverrideWithoutMutatingBase (7.00s)
=== RUN   TestListModelDetails_OverrideIndexIsDeduplicated
    providers_test.go:1150: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_OverrideIndexIsDeduplicated131506862\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_OverrideIndexIsDeduplicated131506862\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_OverrideIndexIsDeduplicated (7.00s)
=== RUN   TestListModelDetails_NoOverridesOmitsNewFields
    providers_test.go:1182: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_NoOverridesOmitsNewFields1043716305\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_NoOverridesOmitsNewFields1043716305\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_NoOverridesOmitsNewFields (7.00s)
=== RUN   TestListModelDetails_PatchToSameValueIsNotMarkedOverridden
    providers_test.go:1197: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_PatchToSameValueIsNotMarkedOverridden2165776567\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_PatchToSameValueIsNotMarkedOverridden2165776567\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_PatchToSameValueIsNotMarkedOverridden (7.00s)
=== RUN   TestListModelDetails_OverrideWithoutBaseCatalogRow
    providers_test.go:1233: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_OverrideWithoutBaseCatalogRow2162828026\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_OverrideWithoutBaseCatalogRow2162828026\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_OverrideWithoutBaseCatalogRow (7.00s)
=== RUN   TestListModelDetails_VirtualKeyScopedOverrideIsInformationalOnly
    providers_test.go:1262: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_VirtualKeyScopedOverrideIsInformationalOnly3536489728\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_VirtualKeyScopedOverrideIsInformationalOnly3536489728\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_VirtualKeyScopedOverrideIsInformationalOnly (7.00s)
=== RUN   TestListModelDetails_AppliesProviderScopedOverride
    providers_test.go:1297: load pricing testdata: failed to load pricing data from URL: failed to parse pricing URL: parse "file://C:\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_AppliesProviderScopedOverride722817774\\001\\pricing.json": invalid port ":\\Users\\MI\\AppData\\Local\\Temp\\TestListModelDetails_AppliesProviderScopedOverride722817774\\001\\pricing.json" after host
--- FAIL: TestListModelDetails_AppliesProviderScopedOverride (7.00s)
=== RUN   TestParseVKValueFromRequest
=== RUN   TestParseVKValueFromRequest/x-bf-vk_header
--- PASS: TestParseVKValueFromRequest/x-bf-vk_header (0.00s)
=== RUN   TestParseVKValueFromRequest/Authorization_Bearer_header
--- PASS: TestParseVKValueFromRequest/Authorization_Bearer_header (0.00s)
=== RUN   TestParseVKValueFromRequest/x-api-key_header
--- PASS: TestParseVKValueFromRequest/x-api-key_header (0.00s)
=== RUN   TestParseVKValueFromRequest/x-goog-api-key_header
--- PASS: TestParseVKValueFromRequest/x-goog-api-key_header (0.00s)
=== RUN   TestParseVKValueFromRequest/no_header_returns_empty_string
--- PASS: TestParseVKValueFromRequest/no_header_returns_empty_string (0.00s)
=== RUN   TestParseVKValueFromRequest/non-VK_Bearer_token_returns_empty_string
--- PASS: TestParseVKValueFromRequest/non-VK_Bearer_token_returns_empty_string (0.00s)
=== RUN   TestParseVKValueFromRequest/x-bf-vk_takes_priority_over_Authorization
--- PASS: TestParseVKValueFromRequest/x-bf-vk_takes_priority_over_Authorization (0.00s)
--- PASS: TestParseVKValueFromRequest (0.00s)
=== RUN   TestListModels_VKFilterHidesBlacklistedModel
--- PASS: TestListModels_VKFilterHidesBlacklistedModel (0.00s)
=== RUN   TestListModels_VKFilterUnionsDuplicateProviderConfigs
--- PASS: TestListModels_VKFilterUnionsDuplicateProviderConfigs (0.00s)
=== RUN   TestListModels_VKFilterRestrictsToAllowedProviderAndModels
--- PASS: TestListModels_VKFilterRestrictsToAllowedProviderAndModels (0.00s)
=== RUN   TestListModels_VKFilterAllowsAllModelsWithWildcard
--- PASS: TestListModels_VKFilterAllowsAllModelsWithWildcard (0.00s)
=== RUN   TestListModels_VKFilterDeniesAllModelsWhenAllowedModelsEmpty
--- PASS: TestListModels_VKFilterDeniesAllModelsWhenAllowedModelsEmpty (0.00s)
=== RUN   TestListModels_VKFilterNoProviderConfigsDeniesAll
--- PASS: TestListModels_VKFilterNoProviderConfigsDeniesAll (0.00s)
=== RUN   TestListModels_VKFilterBlockedExplicitProviderReturnsEmptyResult
--- PASS: TestListModels_VKFilterBlockedExplicitProviderReturnsEmptyResult (0.00s)
=== RUN   TestParseModelListQuery_VKAppliesResolvedAccess
--- PASS: TestParseModelListQuery_VKAppliesResolvedAccess (0.00s)
=== RUN   TestParseModelListQuery_LeavesListingUnfilteredWhenNothingResolved
=== RUN   TestParseModelListQuery_LeavesListingUnfilteredWhenNothingResolved/no_key_presented
--- PASS: TestParseModelListQuery_LeavesListingUnfilteredWhenNothingResolved/no_key_presented (0.00s)
=== RUN   TestParseModelListQuery_LeavesListingUnfilteredWhenNothingResolved/key_resolves_to_nothing
--- PASS: TestParseModelListQuery_LeavesListingUnfilteredWhenNothingResolved/key_resolves_to_nothing (0.00s)
--- PASS: TestParseModelListQuery_LeavesListingUnfilteredWhenNothingResolved (0.00s)
=== RUN   TestListModels_NoVKFilterReturnsAll
--- PASS: TestListModels_NoVKFilterReturnsAll (0.00s)
=== RUN   TestListModels_UsesCatalogAwareAliasMatchingForKeyAllowlist
--- PASS: TestListModels_UsesCatalogAwareAliasMatchingForKeyAllowlist (0.00s)
=== RUN   TestListModels_KeyModelAllowlistIsCaseInsensitive
--- PASS: TestListModels_KeyModelAllowlistIsCaseInsensitive (0.00s)
=== RUN   TestListModels_KeyBlacklistIsCaseInsensitive
--- PASS: TestListModels_KeyBlacklistIsCaseInsensitive (0.00s)
=== RUN   TestRealtimeSessionRoutesOnlyExposeGAClientSecrets
--- PASS: TestRealtimeSessionRoutesOnlyExposeGAClientSecrets (0.00s)
=== RUN   TestResolveRealtimeClientSecretTarget
=== PAUSE TestResolveRealtimeClientSecretTarget
=== RUN   TestResolveRealtimeClientSecretTarget_NormalizesModel
=== PAUSE TestResolveRealtimeClientSecretTarget_NormalizesModel
=== RUN   TestGATranscriptionSessionEndToEndThroughFullNormalizationPath
=== PAUSE TestGATranscriptionSessionEndToEndThroughFullNormalizationPath
=== RUN   TestFullRealtimeSessionWithTranscriptionSiblingEndToEnd
=== PAUSE TestFullRealtimeSessionWithTranscriptionSiblingEndToEnd
=== RUN   TestRewriteGASessionTranscriptionModelAppliesAliasResolution
=== PAUSE TestRewriteGASessionTranscriptionModelAppliesAliasResolution
=== RUN   TestRewriteGASessionTranscriptionModelNoOpForFullRealtimeSession
=== PAUSE TestRewriteGASessionTranscriptionModelNoOpForFullRealtimeSession
=== RUN   TestParseRealtimeEphemeralKeyMapping
=== PAUSE TestParseRealtimeEphemeralKeyMapping
=== RUN   TestParseRealtimeEphemeralKeyMapping_NestedFallback
=== PAUSE TestParseRealtimeEphemeralKeyMapping_NestedFallback
=== RUN   TestParseRealtimeEphemeralKeyMapping_RejectsCompositeShapes
=== PAUSE TestParseRealtimeEphemeralKeyMapping_RejectsCompositeShapes
=== RUN   TestRealtimeMappingVirtualKeyUsesRequestHeader
=== PAUSE TestRealtimeMappingVirtualKeyUsesRequestHeader
=== RUN   TestRealtimeMappingVirtualKeyUsesSettledBearerCredential
=== PAUSE TestRealtimeMappingVirtualKeyUsesSettledBearerCredential
=== RUN   TestRealtimeMappingVirtualKeyFallsBackToContextValue
=== PAUSE TestRealtimeMappingVirtualKeyFallsBackToContextValue
=== RUN   TestReplaceAndCacheRealtimeEphemeralToken
=== PAUSE TestReplaceAndCacheRealtimeEphemeralToken
=== RUN   TestReplaceAndCacheRealtimeEphemeralTokenRejectsExpiredSecret
=== PAUSE TestReplaceAndCacheRealtimeEphemeralTokenRejectsExpiredSecret
=== RUN   TestIsJSONContentType
=== PAUSE TestIsJSONContentType
=== RUN   TestRealtimeClientSecretsEvaluateMintingGovernance_RequiresAccess
=== PAUSE TestRealtimeClientSecretsEvaluateMintingGovernance_RequiresAccess
=== RUN   TestRealtimeClientSecretsEvaluateMintingGovernance_PassesContext
=== PAUSE TestRealtimeClientSecretsEvaluateMintingGovernance_PassesContext
=== RUN   TestRealtimeClientSecretsEvaluateMintingGovernance_ContinuesWithoutGovernance
=== PAUSE TestRealtimeClientSecretsEvaluateMintingGovernance_ContinuesWithoutGovernance
=== RUN   TestRealtimeTurnFinalEventUsesTranscriptionCompletionOnlyForTranscriptionSessions
=== PAUSE TestRealtimeTurnFinalEventUsesTranscriptionCompletionOnlyForTranscriptionSessions
=== RUN   TestRealtimeTurnCompletionContentModelsTranscriptAsOutputOnlyForTranscriptionSessions
=== PAUSE TestRealtimeTurnCompletionContentModelsTranscriptAsOutputOnlyForTranscriptionSessions
=== RUN   TestBuildRealtimeTurnPostResponseModelsTranscriptionAsOutput
=== PAUSE TestBuildRealtimeTurnPostResponseModelsTranscriptionAsOutput
=== RUN   TestBuildRealtimeTurnPreRequestIncludesResponseCreateContent
=== PAUSE TestBuildRealtimeTurnPreRequestIncludesResponseCreateContent
=== RUN   TestBuildRealtimeTurnPreRequestCombinesConversationAndResponseInput
=== PAUSE TestBuildRealtimeTurnPreRequestCombinesConversationAndResponseInput
=== RUN   TestBuildRealtimeTurnPreRequestPreservesEmptyResponseTools
=== PAUSE TestBuildRealtimeTurnPreRequestPreservesEmptyResponseTools
=== RUN   TestShouldAccumulateRealtimeOutput
--- PASS: TestShouldAccumulateRealtimeOutput (0.00s)
=== RUN   TestExtractRealtimeTurnSummary
--- PASS: TestExtractRealtimeTurnSummary (0.00s)
=== RUN   TestFinalizedRealtimeInputSummary
--- PASS: TestFinalizedRealtimeInputSummary (0.00s)
=== RUN   TestFinalizedRealtimeToolOutputSummary
--- PASS: TestFinalizedRealtimeToolOutputSummary (0.00s)
=== RUN   TestPendingRealtimeInputUpdate
=== PAUSE TestPendingRealtimeInputUpdate
=== RUN   TestPendingRealtimeToolOutputUpdate
=== PAUSE TestPendingRealtimeToolOutputUpdate
=== RUN   TestRealtimeSessionDedupeNestedRawEvents
=== PAUSE TestRealtimeSessionDedupeNestedRawEvents
=== RUN   TestBuildRealtimeTurnPostResponseUsesFullResponseDonePayload
--- PASS: TestBuildRealtimeTurnPostResponseUsesFullResponseDonePayload (0.01s)
=== RUN   TestBuildRealtimeTurnPostResponseMergesTextAndToolCalls
--- PASS: TestBuildRealtimeTurnPostResponseMergesTextAndToolCalls (0.00s)
=== RUN   TestFinalizeRealtimeTurnHooksWithErrorCompletesActiveHooks
=== PAUSE TestFinalizeRealtimeTurnHooksWithErrorCompletesActiveHooks
=== RUN   TestNewBifrostErrorFromRealtimeErrorCarriesRealtimeMetadata
=== PAUSE TestNewBifrostErrorFromRealtimeErrorCarriesRealtimeMetadata
=== RUN   TestRefuseUnauthenticatedRealtime_AnonymousRefusedWhenEnforced
--- PASS: TestRefuseUnauthenticatedRealtime_AnonymousRefusedWhenEnforced (0.00s)
=== RUN   TestRefuseUnauthenticatedRealtime_AnonymousAllowedWhenNotEnforced
--- PASS: TestRefuseUnauthenticatedRealtime_AnonymousAllowedWhenNotEnforced (0.00s)
=== RUN   TestRefuseUnauthenticatedRealtime_VirtualKeyAllowed
--- PASS: TestRefuseUnauthenticatedRealtime_VirtualKeyAllowed (0.00s)
=== RUN   TestRefuseUnauthenticatedRealtime_MappedEphemeralTokenAllowed
--- PASS: TestRefuseUnauthenticatedRealtime_MappedEphemeralTokenAllowed (0.00s)
=== RUN   TestRefuseUnauthenticatedRealtime_UnmappedBifrostPrefixedTokenAllowed
--- PASS: TestRefuseUnauthenticatedRealtime_UnmappedBifrostPrefixedTokenAllowed (0.00s)
=== RUN   TestRefuseUnauthenticatedRealtime_NonEphemeralBearerRefused
--- PASS: TestRefuseUnauthenticatedRealtime_NonEphemeralBearerRefused (0.00s)
=== RUN   TestRefuseUnauthenticatedRealtime_DirectKeyAllowed
--- PASS: TestRefuseUnauthenticatedRealtime_DirectKeyAllowed (0.00s)
=== RUN   TestWSRealtimeHandleUpgrade_AnonymousRefusedBeforeTargetResolution
--- PASS: TestWSRealtimeHandleUpgrade_AnonymousRefusedBeforeTargetResolution (0.00s)
=== RUN   TestWSRealtimeHandleUpgrade_TargetErrorStaysInBandWhenNotEnforced
--- PASS: TestWSRealtimeHandleUpgrade_TargetErrorStaysInBandWhenNotEnforced (0.00s)
=== RUN   TestRefuseUnresolvedRealtimeCredential_ForgedVirtualKeyRefused
--- PASS: TestRefuseUnresolvedRealtimeCredential_ForgedVirtualKeyRefused (0.00s)
=== RUN   TestRefuseUnresolvedRealtimeCredential_ResolvedVirtualKeyAdmitted
--- PASS: TestRefuseUnresolvedRealtimeCredential_ResolvedVirtualKeyAdmitted (0.00s)
=== RUN   TestRefuseUnresolvedRealtimeCredential_EphemeralTokenExempt
--- PASS: TestRefuseUnresolvedRealtimeCredential_EphemeralTokenExempt (0.00s)
=== RUN   TestRefuseUnresolvedRealtimeCredential_MappedEphemeralTokenNotExempt
--- PASS: TestRefuseUnresolvedRealtimeCredential_MappedEphemeralTokenNotExempt (0.00s)
=== RUN   TestRefuseUnresolvedRealtimeCredential_ProviderTokenOnlyMappingExempt
--- PASS: TestRefuseUnresolvedRealtimeCredential_ProviderTokenOnlyMappingExempt (0.00s)
=== RUN   TestRefuseUnresolvedRealtimeCredential_DirectKeyAdmitted
--- PASS: TestRefuseUnresolvedRealtimeCredential_DirectKeyAdmitted (0.00s)
=== RUN   TestRefuseUnresolvedRealtimeCredential_NotEnforcedAllowsEverything
--- PASS: TestRefuseUnresolvedRealtimeCredential_NotEnforcedAllowsEverything (0.00s)
=== RUN   TestSessionLoginInvalidPayloadDoesNotExposeDecoderDetails
--- PASS: TestSessionLoginInvalidPayloadDoesNotExposeDecoderDetails (0.00s)
=== RUN   TestPrepareRequestInvalidPayloadDoesNotExposeDecoderDetails
--- PASS: TestPrepareRequestInvalidPayloadDoesNotExposeDecoderDetails (0.00s)
=== RUN   TestComplexityAnalyzerConfigGetDegradesOnUnreadableConfig
    routing_test.go:80: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/routing_test.go:80
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestComplexityAnalyzerConfigGetDegradesOnUnreadableConfig
--- FAIL: TestComplexityAnalyzerConfigGetDegradesOnUnreadableConfig (0.00s)
=== RUN   TestComplexityAnalyzerConfigGetStillFailsWhenStoreUnreachable
--- PASS: TestComplexityAnalyzerConfigGetStillFailsWhenStoreUnreachable (0.00s)
=== RUN   TestComplexityAnalyzerConfigGetReturnsDefaultsWhenUnset
    routing_test.go:124: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/routing_test.go:124
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestComplexityAnalyzerConfigGetReturnsDefaultsWhenUnset
--- FAIL: TestComplexityAnalyzerConfigGetReturnsDefaultsWhenUnset (0.00s)
=== RUN   TestRetryComplexitySemanticWarmup
=== RUN   TestRetryComplexitySemanticWarmup/accepts_a_failed_warmup
    routing_test.go:154: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/routing_test.go:154
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestRetryComplexitySemanticWarmup/accepts_a_failed_warmup
--- FAIL: TestRetryComplexitySemanticWarmup/accepts_a_failed_warmup (0.00s)
=== RUN   TestRetryComplexitySemanticWarmup/rejects_a_non-failed_warmup
    routing_test.go:172: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/routing_test.go:172
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestRetryComplexitySemanticWarmup/rejects_a_non-failed_warmup
--- FAIL: TestRetryComplexitySemanticWarmup/rejects_a_non-failed_warmup (0.00s)
--- FAIL: TestRetryComplexitySemanticWarmup (0.00s)
=== RUN   TestComplexityAnalyzerConfigPutPersistsAndReloads
    routing_test.go:185: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/routing_test.go:185
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestComplexityAnalyzerConfigPutPersistsAndReloads
--- FAIL: TestComplexityAnalyzerConfigPutPersistsAndReloads (0.00s)
=== RUN   TestComplexityAnalyzerConfigPutRejectsInvalidPayloads
    routing_test.go:233: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/routing_test.go:233
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestComplexityAnalyzerConfigPutRejectsInvalidPayloads
--- FAIL: TestComplexityAnalyzerConfigPutRejectsInvalidPayloads (0.00s)
=== RUN   TestComplexityAnalyzerConfigResetPersistsDefaultsAndReloads
    routing_test.go:285: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/routing_test.go:285
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestComplexityAnalyzerConfigResetPersistsDefaultsAndReloads
--- FAIL: TestComplexityAnalyzerConfigResetPersistsDefaultsAndReloads (0.00s)
=== RUN   TestComplexityAnalyzerConfigResetReportsReloadFailure
    routing_test.go:393: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/pricing_override_test.go:86
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/routing_test.go:393
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestComplexityAnalyzerConfigResetReportsReloadFailure
--- FAIL: TestComplexityAnalyzerConfigResetReportsReloadFailure (0.00s)
=== RUN   TestRoutingRoutesServeCanonicalAndLegacyPaths
--- PASS: TestRoutingRoutesServeCanonicalAndLegacyPaths (0.00s)
=== RUN   TestCleanupOrphanSkillFilesDeletesDBFallbackBlobs
    skills_cleanup_test.go:48: new config store: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
--- FAIL: TestCleanupOrphanSkillFilesDeletesDBFallbackBlobs (0.00s)
=== RUN   TestCleanupOrphanSkillFilesDeletesOnlyUnreferencedUploadObjects
    skills_cleanup_test.go:102: new config store: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
--- FAIL: TestCleanupOrphanSkillFilesDeletesOnlyUnreferencedUploadObjects (0.00s)
=== RUN   TestSkillsServingGenericFileDownloadDecodesEncodedPathParams
    skills_serving_test.go:22: new config store: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
--- FAIL: TestSkillsServingGenericFileDownloadDecodesEncodedPathParams (0.00s)
=== RUN   TestClaudeMarketplaceGitRepoContainsMarketplaceAndCloneablePlugin
    skills_serving_test.go:88: new config store: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
--- FAIL: TestClaudeMarketplaceGitRepoContainsMarketplaceAndCloneablePlugin (0.00s)
=== RUN   TestPrepareSpeechRequestVoiceValidation
=== RUN   TestPrepareSpeechRequestVoiceValidation/elevenlabs_alias_without_voice_is_deferred_to_provider
--- PASS: TestPrepareSpeechRequestVoiceValidation/elevenlabs_alias_without_voice_is_deferred_to_provider (0.00s)
=== RUN   TestPrepareSpeechRequestVoiceValidation/elevenlabs_literal_sound_model_without_voice_is_allowed
--- PASS: TestPrepareSpeechRequestVoiceValidation/elevenlabs_literal_sound_model_without_voice_is_allowed (0.00s)
=== RUN   TestPrepareSpeechRequestVoiceValidation/elevenlabs_tts_without_voice_is_deferred_to_provider_(no_transport_400)
--- PASS: TestPrepareSpeechRequestVoiceValidation/elevenlabs_tts_without_voice_is_deferred_to_provider_(no_transport_400) (0.00s)
=== RUN   TestPrepareSpeechRequestVoiceValidation/non-elevenlabs_provider_without_voice_still_errors_at_the_transport
--- PASS: TestPrepareSpeechRequestVoiceValidation/non-elevenlabs_provider_without_voice_still_errors_at_the_transport (0.00s)
--- PASS: TestPrepareSpeechRequestVoiceValidation (0.00s)
=== RUN   TestSSEStreamReaderNoEventBatching
--- PASS: TestSSEStreamReaderNoEventBatching (0.00s)
=== RUN   TestProbeChatRequestStopParse
--- PASS: TestProbeChatRequestStopParse (0.04s)
=== RUN   TestStreamingResponseSkipsDoneAfterErrorChunk
--- PASS: TestStreamingResponseSkipsDoneAfterErrorChunk (0.21s)
=== RUN   TestStreamingResponseSendsDoneOnCleanStream
--- PASS: TestStreamingResponseSendsDoneOnCleanStream (0.20s)
=== RUN   TestSendJSON_StandardBytes
--- PASS: TestSendJSON_StandardBytes (0.00s)
=== RUN   TestSendJSON_Deterministic
--- PASS: TestSendJSON_Deterministic (0.00s)
=== RUN   TestSendJSONWithStatus_StandardBytes
--- PASS: TestSendJSONWithStatus_StandardBytes (0.00s)
=== RUN   TestSendJSON_MarshalError
--- PASS: TestSendJSON_MarshalError (0.00s)
=== RUN   TestSendJSONWithStatus_MarshalError
--- PASS: TestSendJSONWithStatus_MarshalError (0.00s)
=== RUN   TestSendBifrostError_BodylessStatusBecomes502
--- PASS: TestSendBifrostError_BodylessStatusBecomes502 (0.00s)
=== RUN   TestSendBifrostError_NormalStatusPreserved
--- PASS: TestSendBifrostError_NormalStatusPreserved (0.00s)
=== RUN   TestWebhookRoutesRegister
--- PASS: TestWebhookRoutesRegister (0.00s)
=== RUN   TestWebhookHandlerRouteRegistration
    webhooks_test.go:100: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:100
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerRouteRegistration
--- FAIL: TestWebhookHandlerRouteRegistration (0.00s)
=== RUN   TestWebhookHandlerCreateAndGet
    webhooks_test.go:107: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:107
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerCreateAndGet
--- FAIL: TestWebhookHandlerCreateAndGet (0.00s)
=== RUN   TestWebhookHandlerListFilters
    webhooks_test.go:139: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:139
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerListFilters
--- FAIL: TestWebhookHandlerListFilters (0.00s)
=== RUN   TestWebhookHandlerCreateValidation
    webhooks_test.go:199: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:199
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerCreateValidation
--- FAIL: TestWebhookHandlerCreateValidation (0.00s)
=== RUN   TestWebhookHandlerDuplicateName
    webhooks_test.go:216: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:216
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerDuplicateName
--- FAIL: TestWebhookHandlerDuplicateName (0.00s)
=== RUN   TestWebhookHandlerUpdate
    webhooks_test.go:225: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:225
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerUpdate
--- FAIL: TestWebhookHandlerUpdate (0.00s)
=== RUN   TestWebhookHandlerDelete
    webhooks_test.go:246: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:246
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerDelete
--- FAIL: TestWebhookHandlerDelete (0.00s)
=== RUN   TestWebhookHandlerRotateSecret
    webhooks_test.go:262: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:262
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerRotateSecret
--- FAIL: TestWebhookHandlerRotateSecret (0.00s)
=== RUN   TestWebhookHandlerListDeliveries
    webhooks_test.go:293: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:293
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerListDeliveries
--- FAIL: TestWebhookHandlerListDeliveries (0.00s)
=== RUN   TestWebhookHandlerRedeliver
    webhooks_test.go:317: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:317
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerRedeliver
--- FAIL: TestWebhookHandlerRedeliver (0.00s)
=== RUN   TestWebhookHandlerTestDelivery
    webhooks_test.go:349: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:349
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerTestDelivery
--- FAIL: TestWebhookHandlerTestDelivery (0.00s)
=== RUN   TestWebhookHandlerStoreUnavailable
--- PASS: TestWebhookHandlerStoreUnavailable (0.00s)
=== RUN   TestWebhookHandlerPerEndpointTuning
    webhooks_test.go:406: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:406
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerPerEndpointTuning
--- FAIL: TestWebhookHandlerPerEndpointTuning (0.00s)
=== RUN   TestWebhookHandlerHeaders
    webhooks_test.go:430: 
        	Error Trace:	C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:53
        	            				C:/Users/MI/aigw-development/worktrees/bifrost-cookie67-baseline/transports/bifrost-http/handlers/webhooks_test.go:430
        	Error:      	Received unexpected error:
        	            	Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
        	Test:       	TestWebhookHandlerHeaders
--- FAIL: TestWebhookHandlerHeaders (0.00s)
=== RUN   TestExtractRealtimeTokenFromAuth
=== PAUSE TestExtractRealtimeTokenFromAuth
=== RUN   TestResolveRealtimeSDPTarget_BaseRouteRequiresProviderPrefix
--- PASS: TestResolveRealtimeSDPTarget_BaseRouteRequiresProviderPrefix (0.00s)
=== RUN   TestResolveRealtimeSDPTarget_BaseRouteNormalizesModel
--- PASS: TestResolveRealtimeSDPTarget_BaseRouteNormalizesModel (0.00s)
=== RUN   TestResolveRealtimeSDPTarget_OpenAIRouteDefaultsProvider
--- PASS: TestResolveRealtimeSDPTarget_OpenAIRouteDefaultsProvider (0.00s)
=== RUN   TestParseCallsWebRTCRequest_RawSDPKeepsGARoute
--- PASS: TestParseCallsWebRTCRequest_RawSDPKeepsGARoute (0.00s)
=== RUN   TestResolveRealtimeSDPTargetDedicatedTranscription
=== PAUSE TestResolveRealtimeSDPTargetDedicatedTranscription
=== RUN   TestResolveRealtimeSDPTargetRootModelPreservesNormalRealtime
=== PAUSE TestResolveRealtimeSDPTargetRootModelPreservesNormalRealtime
=== RUN   TestPinRealtimeSDPTranscriptionModelPreservesSession
=== PAUSE TestPinRealtimeSDPTranscriptionModelPreservesSession
=== RUN   TestNewRealtimeRelayContextCarriesRequestGrant
--- PASS: TestNewRealtimeRelayContextCarriesRequestGrant (0.00s)
=== RUN   TestNewRealtimeRelayContextCopiesValuesWithoutRequestCancellation
--- PASS: TestNewRealtimeRelayContextCopiesValuesWithoutRequestCancellation (0.00s)
=== RUN   TestParseRealtimeEventPreservesExtraParams
--- PASS: TestParseRealtimeEventPreservesExtraParams (0.01s)
=== RUN   TestExtractRealtimeBearerToken
--- PASS: TestExtractRealtimeBearerToken (0.00s)
=== RUN   TestLookupRealtimeEphemeralKeyMappingKeepsEntryUntilTTLExpiry
=== PAUSE TestLookupRealtimeEphemeralKeyMappingKeepsEntryUntilTTLExpiry
=== RUN   TestCacheRealtimeEphemeralKeyMappingPreservesVirtualKeyAcrossKVReplication
=== PAUSE TestCacheRealtimeEphemeralKeyMappingPreservesVirtualKeyAcrossKVReplication
=== RUN   TestLookupRealtimeEphemeralKeyMapping_BackwardsCompatibleStringValue
=== PAUSE TestLookupRealtimeEphemeralKeyMapping_BackwardsCompatibleStringValue
=== RUN   TestWebRTCRealtimeRelayCloseFinalizesActiveTurnHooks
=== PAUSE TestWebRTCRealtimeRelayCloseFinalizesActiveTurnHooks
=== RUN   TestResolveRealtimeWebRTCKeys_UnmappedEphemeralTokenStaysAnonymous
=== PAUSE TestResolveRealtimeWebRTCKeys_UnmappedEphemeralTokenStaysAnonymous
=== RUN   TestApplyRealtimeEphemeralKeyMapping_RestoresVirtualKeyAndKeyID
=== PAUSE TestApplyRealtimeEphemeralKeyMapping_RestoresVirtualKeyAndKeyID
=== RUN   TestBroadcastNotificationTargetsMatchingRole
--- PASS: TestBroadcastNotificationTargetsMatchingRole (0.21s)
=== RUN   TestWebSocketWriteAfterHandlerReturn
--- PASS: TestWebSocketWriteAfterHandlerReturn (0.25s)
=== RUN   TestWebSocketBroadcastDuringDisconnectStorm
--- PASS: TestWebSocketBroadcastDuringDisconnectStorm (0.09s)
=== RUN   TestWebSocketStopClosesClients
--- PASS: TestWebSocketStopClosesClients (0.00s)
=== RUN   TestMatchesWildcardPattern
=== RUN   TestMatchesWildcardPattern/subdomain_match
--- PASS: TestMatchesWildcardPattern/subdomain_match (0.00s)
=== RUN   TestMatchesWildcardPattern/subdomain_no_match
--- PASS: TestMatchesWildcardPattern/subdomain_no_match (0.00s)
=== RUN   TestMatchesWildcardPattern/scheme-less_no_match_with_scheme_prefix
--- PASS: TestMatchesWildcardPattern/scheme-less_no_match_with_scheme_prefix (0.00s)
=== RUN   TestMatchesWildcardPattern/scheme-less_exact_subdomain
--- PASS: TestMatchesWildcardPattern/scheme-less_exact_subdomain (0.00s)
=== RUN   TestMatchesWildcardPattern/no_nested_subdomain
--- PASS: TestMatchesWildcardPattern/no_nested_subdomain (0.00s)
=== RUN   TestMatchesWildcardPattern/empty_subdomain_no_match
--- PASS: TestMatchesWildcardPattern/empty_subdomain_no_match (0.00s)
=== RUN   TestMatchesWildcardPattern/no_wildcard_in_pattern
--- PASS: TestMatchesWildcardPattern/no_wildcard_in_pattern (0.00s)
=== RUN   TestMatchesWildcardPattern/http_vs_https
--- PASS: TestMatchesWildcardPattern/http_vs_https (0.00s)
=== RUN   TestMatchesWildcardPattern/http_wildcard
--- PASS: TestMatchesWildcardPattern/http_wildcard (0.00s)
=== RUN   TestMatchesWildcardPattern/multiple_wildcards
--- PASS: TestMatchesWildcardPattern/multiple_wildcards (0.00s)
=== RUN   TestMatchesWildcardPattern/empty_origin
--- PASS: TestMatchesWildcardPattern/empty_origin (0.00s)
=== RUN   TestMatchesWildcardPattern/empty_pattern
--- PASS: TestMatchesWildcardPattern/empty_pattern (0.00s)
=== RUN   TestMatchesWildcardPattern/both_empty
--- PASS: TestMatchesWildcardPattern/both_empty (0.00s)
=== RUN   TestMatchesWildcardPattern/no_slash_in_wildcard
--- PASS: TestMatchesWildcardPattern/no_slash_in_wildcard (0.00s)
--- PASS: TestMatchesWildcardPattern (0.00s)
=== RUN   TestMatchesWildcardPattern_CacheConsistency
--- PASS: TestMatchesWildcardPattern_CacheConsistency (0.00s)
=== RUN   TestIsOriginAllowed
=== RUN   TestIsOriginAllowed/localhost_http
--- PASS: TestIsOriginAllowed/localhost_http (0.00s)
=== RUN   TestIsOriginAllowed/localhost_https
--- PASS: TestIsOriginAllowed/localhost_https (0.00s)
=== RUN   TestIsOriginAllowed/127.0.0.1
--- PASS: TestIsOriginAllowed/127.0.0.1 (0.00s)
=== RUN   TestIsOriginAllowed/exact_match
--- PASS: TestIsOriginAllowed/exact_match (0.00s)
=== RUN   TestIsOriginAllowed/exact_no_match
--- PASS: TestIsOriginAllowed/exact_no_match (0.00s)
=== RUN   TestIsOriginAllowed/star_allows_all
--- PASS: TestIsOriginAllowed/star_allows_all (0.00s)
=== RUN   TestIsOriginAllowed/wildcard_match
--- PASS: TestIsOriginAllowed/wildcard_match (0.00s)
=== RUN   TestIsOriginAllowed/wildcard_no_match
--- PASS: TestIsOriginAllowed/wildcard_no_match (0.00s)
=== RUN   TestIsOriginAllowed/exact_before_wildcard
--- PASS: TestIsOriginAllowed/exact_before_wildcard (0.00s)
=== RUN   TestIsOriginAllowed/wildcard_after_exact_miss
--- PASS: TestIsOriginAllowed/wildcard_after_exact_miss (0.00s)
=== RUN   TestIsOriginAllowed/empty_origins
--- PASS: TestIsOriginAllowed/empty_origins (0.00s)
=== RUN   TestIsOriginAllowed/empty_origins_slice
--- PASS: TestIsOriginAllowed/empty_origins_slice (0.00s)
--- PASS: TestIsOriginAllowed (0.00s)
=== RUN   TestDiscoverRealtimeTranscriptionModel
=== PAUSE TestDiscoverRealtimeTranscriptionModel
=== RUN   TestPinRealtimeTranscriptionModel
=== PAUSE TestPinRealtimeTranscriptionModel
=== RUN   TestSnapshotRealtimeMiddlewareValuesWithContext
=== PAUSE TestSnapshotRealtimeMiddlewareValuesWithContext
=== RUN   TestBufferRealtimeTranscriptionBootstrapPreservesFrames
--- PASS: TestBufferRealtimeTranscriptionBootstrapPreservesFrames (0.00s)
=== RUN   TestRealtimeStopHeartbeatWaitsForPingGoroutine
--- PASS: TestRealtimeStopHeartbeatWaitsForPingGoroutine (0.05s)
=== RUN   TestRealtimeStopHeartbeatWithoutStart
--- PASS: TestRealtimeStopHeartbeatWithoutStart (0.00s)
=== RUN   TestRealtimeStopHeartbeatIsIdempotent
--- PASS: TestRealtimeStopHeartbeatIsIdempotent (0.12s)
=== RUN   TestNativeWSUpstreamProducesLLMCallSpan
--- PASS: TestNativeWSUpstreamProducesLLMCallSpan (0.10s)
=== RUN   TestNativeWSUpstreamFailedResponseProducesErroredLLMCallSpan
--- PASS: TestNativeWSUpstreamFailedResponseProducesErroredLLMCallSpan (0.00s)
=== RUN   TestNativeWSUpstreamTimeoutProducesErroredLLMCallSpan
--- PASS: TestNativeWSUpstreamTimeoutProducesErroredLLMCallSpan (1.00s)
=== RUN   TestChatGPTWSMultiTurnRawAndQuota
=== RUN   TestChatGPTWSMultiTurnRawAndQuota/false
--- PASS: TestChatGPTWSMultiTurnRawAndQuota/false (0.01s)
=== RUN   TestChatGPTWSMultiTurnRawAndQuota/true
--- PASS: TestChatGPTWSMultiTurnRawAndQuota/true (0.00s)
--- PASS: TestChatGPTWSMultiTurnRawAndQuota (0.01s)
=== RUN   TestChatGPTWSTargetRejectsPlaintextRemote
--- PASS: TestChatGPTWSTargetRejectsPlaintextRemote (0.00s)
=== RUN   TestChatGPTWSRouteRequiresUpgradeAndCredentials
--- PASS: TestChatGPTWSRouteRequiresUpgradeAndCredentials (0.00s)
=== RUN   TestChatGPTWSDisconnectDuringTurnFinalizesHooks
--- PASS: TestChatGPTWSDisconnectDuringTurnFinalizesHooks (0.00s)
=== RUN   TestResolveWSStreamIdleTimeoutUsesProviderOverride
--- PASS: TestResolveWSStreamIdleTimeoutUsesProviderOverride (0.00s)
=== RUN   TestResolveWSStreamIdleTimeoutFallsBackToDefault
--- PASS: TestResolveWSStreamIdleTimeoutFallsBackToDefault (0.00s)
=== RUN   TestIsWSReadTimeout
--- PASS: TestIsWSReadTimeout (0.00s)
=== RUN   TestNewBifrostError
--- PASS: TestNewBifrostError (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_BaggageSessionIDSetsGrouping
--- PASS: TestCreateBifrostContextFromAuth_BaggageSessionIDSetsGrouping (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_SettlesIdentity
=== RUN   TestCreateBifrostContextFromAuth_SettlesIdentity/virtual_key_header
--- PASS: TestCreateBifrostContextFromAuth_SettlesIdentity/virtual_key_header (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_SettlesIdentity/bearer_virtual_key
--- PASS: TestCreateBifrostContextFromAuth_SettlesIdentity/bearer_virtual_key (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_SettlesIdentity/provider_key_bearer
--- PASS: TestCreateBifrostContextFromAuth_SettlesIdentity/provider_key_bearer (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_SettlesIdentity/no_credential
--- PASS: TestCreateBifrostContextFromAuth_SettlesIdentity/no_credential (0.00s)
--- PASS: TestCreateBifrostContextFromAuth_SettlesIdentity (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_EmptyBaggageSessionIDIgnored
--- PASS: TestCreateBifrostContextFromAuth_EmptyBaggageSessionIDIgnored (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_ForwardsPrefixedHeaders
--- PASS: TestCreateBifrostContextFromAuth_ForwardsPrefixedHeaders (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_AppliesHeaderFilterAndDirectAllowlist
--- PASS: TestCreateBifrostContextFromAuth_AppliesHeaderFilterAndDirectAllowlist (0.00s)
=== RUN   TestCaptureAuthHeaders_PreservesDuplicateHeaderValues
--- PASS: TestCaptureAuthHeaders_PreservesDuplicateHeaderValues (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_PreservesMultipleForwardedHeaderValues
--- PASS: TestCreateBifrostContextFromAuth_PreservesMultipleForwardedHeaderValues (0.00s)
=== RUN   TestCreateBifrostContextFromAuth_BlocksWebSocketHandshakeForwardedHeaders
--- PASS: TestCreateBifrostContextFromAuth_BlocksWebSocketHandshakeForwardedHeaders (0.00s)
=== RUN   TestMergeWebSocketHeaders_ForwardedHeadersOverrideProviderHeadersAndPreserveValues
--- PASS: TestMergeWebSocketHeaders_ForwardedHeadersOverrideProviderHeadersAndPreserveValues (0.00s)
=== RUN   TestHasWebSocketForwardedHeaders
--- PASS: TestHasWebSocketForwardedHeaders (0.00s)
=== RUN   TestSignedWSTicketValidatesAcrossStores
--- PASS: TestSignedWSTicketValidatesAcrossStores (0.00s)
=== RUN   TestSignedWSTicketRejectsWrongKey
--- PASS: TestSignedWSTicketRejectsWrongKey (0.00s)
=== RUN   TestSignedWSTicketRejectsExpiredTicket
--- PASS: TestSignedWSTicketRejectsExpiredTicket (0.00s)
=== RUN   TestSignedWSTicketRejectsMalformedTicket
--- PASS: TestSignedWSTicketRejectsMalformedTicket (0.00s)
=== RUN   TestSignedWSTicketRejectsTamperedTicket
--- PASS: TestSignedWSTicketRejectsTamperedTicket (0.00s)
=== RUN   TestLegacyWSTicketRemainsSingleUse
--- PASS: TestLegacyWSTicketRemainsSingleUse (0.00s)
=== RUN   TestSignedWSTicketDoesNotExposeSessionToken
--- PASS: TestSignedWSTicketDoesNotExposeSessionToken (0.00s)
=== RUN   TestSignedWSTicketPayloadEncryptionRejectsShortPayload
--- PASS: TestSignedWSTicketPayloadEncryptionRejectsShortPayload (0.00s)
=== RUN   TestSignedWSTicketNonceIsHex
--- PASS: TestSignedWSTicketNonceIsHex (0.00s)
=== CONT  TestResolveRealtimeClientSecretTarget
=== CONT  TestBuildRealtimeTurnPreRequestIncludesResponseCreateContent
=== CONT  TestRealtimeClientSecretsEvaluateMintingGovernance_PassesContext
=== CONT  TestParseRealtimeEphemeralKeyMapping_RejectsCompositeShapes
=== CONT  TestIsJSONContentType
--- PASS: TestRealtimeClientSecretsEvaluateMintingGovernance_PassesContext (0.00s)
=== CONT  TestParseRealtimeEphemeralKeyMapping_NestedFallback
=== CONT  TestReplaceAndCacheRealtimeEphemeralTokenRejectsExpiredSecret
=== CONT  TestRealtimeMappingVirtualKeyFallsBackToContextValue
--- PASS: TestIsJSONContentType (0.00s)
=== CONT  TestRealtimeMappingVirtualKeyUsesSettledBearerCredential
--- PASS: TestRealtimeMappingVirtualKeyFallsBackToContextValue (0.00s)
=== CONT  TestResolveRealtimeClientSecretTarget_NormalizesModel
=== CONT  TestRealtimeMappingVirtualKeyUsesRequestHeader
=== RUN   TestResolveRealtimeClientSecretTarget_NormalizesModel/session.model_provider_prefix_stripped
=== CONT  TestReplaceAndCacheRealtimeEphemeralToken
=== CONT  TestSnapshotRealtimeMiddlewareValuesWithContext
=== CONT  TestRealtimeTurnFinalEventUsesTranscriptionCompletionOnlyForTranscriptionSessions
=== CONT  TestRealtimeClientSecretsEvaluateMintingGovernance_ContinuesWithoutGovernance
=== CONT  TestPinRealtimeTranscriptionModel
=== CONT  TestRealtimeTurnCompletionContentModelsTranscriptAsOutputOnlyForTranscriptionSessions
=== CONT  TestParseRealtimeEphemeralKeyMapping
=== CONT  TestRewriteGASessionTranscriptionModelNoOpForFullRealtimeSession
=== CONT  TestResolveRealtimeWebRTCKeys_UnmappedEphemeralTokenStaysAnonymous
=== CONT  TestRewriteGASessionTranscriptionModelAppliesAliasResolution
=== CONT  TestWebRTCRealtimeRelayCloseFinalizesActiveTurnHooks
=== CONT  TestLookupRealtimeEphemeralKeyMapping_BackwardsCompatibleStringValue
=== CONT  TestBuildRealtimeTurnPostResponseModelsTranscriptionAsOutput
=== RUN   TestResolveRealtimeClientSecretTarget/base_route_with_session_model
=== CONT  TestCacheRealtimeEphemeralKeyMappingPreservesVirtualKeyAcrossKVReplication
=== CONT  TestLookupRealtimeEphemeralKeyMappingKeepsEntryUntilTTLExpiry
=== PAUSE TestResolveRealtimeClientSecretTarget/base_route_with_session_model
=== CONT  TestFullRealtimeSessionWithTranscriptionSiblingEndToEnd
=== CONT  TestPendingRealtimeInputUpdate
=== CONT  TestRealtimeSessionDedupeNestedRawEvents
=== CONT  TestPendingRealtimeToolOutputUpdate
=== CONT  TestBuildRealtimeTurnPreRequestPreservesEmptyResponseTools
=== CONT  TestBuildRealtimeTurnPreRequestCombinesConversationAndResponseInput
=== CONT  TestFinalizeRealtimeTurnHooksWithErrorCompletesActiveHooks
=== CONT  TestPinRealtimeSDPTranscriptionModelPreservesSession
=== CONT  TestDiscoverRealtimeTranscriptionModel
=== CONT  TestRealtimeClientSecretsEvaluateMintingGovernance_RequiresAccess
=== RUN   TestPinRealtimeTranscriptionModel/transcription_session_uses_pinned_alias-resolved_model
=== CONT  TestApplyRealtimeEphemeralKeyMapping_RestoresVirtualKeyAndKeyID
=== CONT  TestNewBifrostErrorFromRealtimeErrorCarriesRealtimeMetadata
=== CONT  TestGATranscriptionSessionEndToEndThroughFullNormalizationPath
=== CONT  TestResolveRealtimeSDPTargetDedicatedTranscription
=== CONT  TestResolveRealtimeSDPTargetRootModelPreservesNormalRealtime
=== RUN   TestResolveRealtimeClientSecretTarget/openai_alias_uses_bare_model
--- PASS: TestParseRealtimeEphemeralKeyMapping_RejectsCompositeShapes (0.00s)
=== PAUSE TestResolveRealtimeClientSecretTarget/openai_alias_uses_bare_model
--- PASS: TestRealtimeMappingVirtualKeyUsesSettledBearerCredential (0.00s)
=== CONT  TestExtractRealtimeTokenFromAuth
--- PASS: TestParseRealtimeEphemeralKeyMapping_NestedFallback (0.00s)
--- PASS: TestReplaceAndCacheRealtimeEphemeralTokenRejectsExpiredSecret (0.00s)
--- PASS: TestRealtimeMappingVirtualKeyUsesRequestHeader (0.00s)
--- PASS: TestRealtimeClientSecretsEvaluateMintingGovernance_ContinuesWithoutGovernance (0.00s)
=== RUN   TestExtractRealtimeTokenFromAuth/google_api_key
--- PASS: TestResolveRealtimeSDPTargetDedicatedTranscription (0.00s)
=== PAUSE TestExtractRealtimeTokenFromAuth/google_api_key
--- PASS: TestSnapshotRealtimeMiddlewareValuesWithContext (0.00s)
--- PASS: TestRealtimeTurnFinalEventUsesTranscriptionCompletionOnlyForTranscriptionSessions (0.00s)
=== PAUSE TestResolveRealtimeClientSecretTarget_NormalizesModel/session.model_provider_prefix_stripped
=== RUN   TestExtractRealtimeTokenFromAuth/authorization_wins
=== PAUSE TestExtractRealtimeTokenFromAuth/authorization_wins
=== RUN   TestDiscoverRealtimeTranscriptionModel/nested_transcription_model
=== RUN   TestResolveRealtimeClientSecretTarget_NormalizesModel/bare_model_unchanged_on_alias_route
=== PAUSE TestResolveRealtimeClientSecretTarget_NormalizesModel/bare_model_unchanged_on_alias_route
=== RUN   TestExtractRealtimeTokenFromAuth/nil
=== CONT  TestResolveRealtimeClientSecretTarget_NormalizesModel/session.model_provider_prefix_stripped
--- PASS: TestRealtimeTurnCompletionContentModelsTranscriptAsOutputOnlyForTranscriptionSessions (0.00s)
=== CONT  TestResolveRealtimeClientSecretTarget_NormalizesModel/bare_model_unchanged_on_alias_route
--- PASS: TestRewriteGASessionTranscriptionModelNoOpForFullRealtimeSession (0.00s)
--- PASS: TestParseRealtimeEphemeralKeyMapping (0.00s)
--- PASS: TestLookupRealtimeEphemeralKeyMapping_BackwardsCompatibleStringValue (0.00s)
=== RUN   TestResolveRealtimeClientSecretTarget/base_route_rejects_bare_model
--- PASS: TestWebRTCRealtimeRelayCloseFinalizesActiveTurnHooks (0.00s)
=== PAUSE TestPinRealtimeTranscriptionModel/transcription_session_uses_pinned_alias-resolved_model
=== PAUSE TestDiscoverRealtimeTranscriptionModel/nested_transcription_model
=== RUN   TestPinRealtimeTranscriptionModel/normal_realtime_keeps_input_transcription_model
=== PAUSE TestPinRealtimeTranscriptionModel/normal_realtime_keeps_input_transcription_model
=== CONT  TestPinRealtimeTranscriptionModel/transcription_session_uses_pinned_alias-resolved_model
=== RUN   TestDiscoverRealtimeTranscriptionModel/wrong_event_type
=== CONT  TestPinRealtimeTranscriptionModel/normal_realtime_keeps_input_transcription_model
=== PAUSE TestDiscoverRealtimeTranscriptionModel/wrong_event_type
=== RUN   TestDiscoverRealtimeTranscriptionModel/top-level_model
--- PASS: TestResolveRealtimeClientSecretTarget_NormalizesModel/bare_model_unchanged_on_alias_route (0.00s)
=== PAUSE TestDiscoverRealtimeTranscriptionModel/top-level_model
--- PASS: TestResolveRealtimeClientSecretTarget_NormalizesModel/session.model_provider_prefix_stripped (0.00s)
--- PASS: TestResolveRealtimeWebRTCKeys_UnmappedEphemeralTokenStaysAnonymous (0.00s)
=== RUN   TestDiscoverRealtimeTranscriptionModel/invalid_JSON
=== PAUSE TestDiscoverRealtimeTranscriptionModel/invalid_JSON
=== PAUSE TestExtractRealtimeTokenFromAuth/nil
=== CONT  TestDiscoverRealtimeTranscriptionModel/wrong_event_type
=== RUN   TestExtractRealtimeTokenFromAuth/authorization_bearer
--- PASS: TestReplaceAndCacheRealtimeEphemeralToken (0.00s)
=== PAUSE TestExtractRealtimeTokenFromAuth/authorization_bearer
--- PASS: TestDiscoverRealtimeTranscriptionModel/wrong_event_type (0.00s)
--- PASS: TestRewriteGASessionTranscriptionModelAppliesAliasResolution (0.00s)
--- PASS: TestBuildRealtimeTurnPostResponseModelsTranscriptionAsOutput (0.00s)
=== CONT  TestDiscoverRealtimeTranscriptionModel/nested_transcription_model
=== PAUSE TestResolveRealtimeClientSecretTarget/base_route_rejects_bare_model
=== CONT  TestDiscoverRealtimeTranscriptionModel/invalid_JSON
=== RUN   TestResolveRealtimeClientSecretTarget/missing_model
=== PAUSE TestResolveRealtimeClientSecretTarget/missing_model
=== CONT  TestDiscoverRealtimeTranscriptionModel/top-level_model
=== RUN   TestExtractRealtimeTokenFromAuth/virtual_key_header
=== PAUSE TestExtractRealtimeTokenFromAuth/virtual_key_header
--- PASS: TestPendingRealtimeInputUpdate (0.00s)
--- PASS: TestPinRealtimeTranscriptionModel/transcription_session_uses_pinned_alias-resolved_model (0.00s)
--- PASS: TestDiscoverRealtimeTranscriptionModel/nested_transcription_model (0.00s)
--- PASS: TestPinRealtimeTranscriptionModel/normal_realtime_keeps_input_transcription_model (0.00s)
--- PASS: TestDiscoverRealtimeTranscriptionModel/invalid_JSON (0.00s)
--- PASS: TestDiscoverRealtimeTranscriptionModel/top-level_model (0.00s)
=== RUN   TestResolveRealtimeClientSecretTarget/GA_transcription_session_resolves_model_from_nested_audio_path
--- PASS: TestRealtimeSessionDedupeNestedRawEvents (0.00s)
--- PASS: TestPendingRealtimeToolOutputUpdate (0.00s)
=== RUN   TestExtractRealtimeTokenFromAuth/api_key_header
=== PAUSE TestResolveRealtimeClientSecretTarget/GA_transcription_session_resolves_model_from_nested_audio_path
=== PAUSE TestExtractRealtimeTokenFromAuth/api_key_header
=== CONT  TestExtractRealtimeTokenFromAuth/google_api_key
=== CONT  TestResolveRealtimeClientSecretTarget/openai_alias_uses_bare_model
--- PASS: TestLookupRealtimeEphemeralKeyMappingKeepsEntryUntilTTLExpiry (0.00s)
--- PASS: TestFinalizeRealtimeTurnHooksWithErrorCompletesActiveHooks (0.00s)
=== CONT  TestResolveRealtimeClientSecretTarget/base_route_with_session_model
=== CONT  TestResolveRealtimeClientSecretTarget/GA_transcription_session_resolves_model_from_nested_audio_path
=== CONT  TestResolveRealtimeClientSecretTarget/missing_model
=== CONT  TestResolveRealtimeClientSecretTarget/base_route_rejects_bare_model
=== CONT  TestExtractRealtimeTokenFromAuth/authorization_bearer
=== CONT  TestExtractRealtimeTokenFromAuth/api_key_header
=== CONT  TestExtractRealtimeTokenFromAuth/virtual_key_header
--- PASS: TestExtractRealtimeTokenFromAuth/google_api_key (0.00s)
--- PASS: TestExtractRealtimeTokenFromAuth/authorization_bearer (0.00s)
--- PASS: TestExtractRealtimeTokenFromAuth/api_key_header (0.00s)
=== CONT  TestExtractRealtimeTokenFromAuth/authorization_wins
--- PASS: TestResolveRealtimeClientSecretTarget/openai_alias_uses_bare_model (0.00s)
=== CONT  TestExtractRealtimeTokenFromAuth/nil
--- PASS: TestBuildRealtimeTurnPreRequestPreservesEmptyResponseTools (0.00s)
--- PASS: TestRealtimeClientSecretsEvaluateMintingGovernance_RequiresAccess (0.00s)
--- PASS: TestExtractRealtimeTokenFromAuth/virtual_key_header (0.00s)
--- PASS: TestResolveRealtimeClientSecretTarget/base_route_rejects_bare_model (0.00s)
--- PASS: TestExtractRealtimeTokenFromAuth/authorization_wins (0.00s)
--- PASS: TestApplyRealtimeEphemeralKeyMapping_RestoresVirtualKeyAndKeyID (0.00s)
--- PASS: TestExtractRealtimeTokenFromAuth/nil (0.00s)
--- PASS: TestResolveRealtimeClientSecretTarget/base_route_with_session_model (0.00s)
--- PASS: TestResolveRealtimeClientSecretTarget/missing_model (0.00s)
--- PASS: TestNewBifrostErrorFromRealtimeErrorCarriesRealtimeMetadata (0.00s)
--- PASS: TestResolveRealtimeClientSecretTarget/GA_transcription_session_resolves_model_from_nested_audio_path (0.00s)
--- PASS: TestFullRealtimeSessionWithTranscriptionSiblingEndToEnd (0.00s)
--- PASS: TestPinRealtimeSDPTranscriptionModelPreservesSession (0.00s)
--- PASS: TestResolveRealtimeSDPTargetRootModelPreservesNormalRealtime (0.00s)
--- PASS: TestGATranscriptionSessionEndToEndThroughFullNormalizationPath (0.00s)
--- PASS: TestResolveRealtimeClientSecretTarget_NormalizesModel (0.00s)
--- PASS: TestPinRealtimeTranscriptionModel (0.00s)
--- PASS: TestCacheRealtimeEphemeralKeyMappingPreservesVirtualKeyAcrossKVReplication (0.00s)
--- PASS: TestDiscoverRealtimeTranscriptionModel (0.00s)
--- PASS: TestResolveRealtimeClientSecretTarget (0.00s)
--- PASS: TestExtractRealtimeTokenFromAuth (0.00s)
--- PASS: TestBuildRealtimeTurnPreRequestCombinesConversationAndResponseInput (0.06s)
--- PASS: TestBuildRealtimeTurnPreRequestIncludesResponseCreateContent (0.06s)
FAIL
FAIL	github.com/maximhq/bifrost/transports/bifrost-http/handlers	117.226s

```
