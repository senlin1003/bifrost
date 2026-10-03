# FWD-67：原生 WS 身份关联

开发者：GPT；2026-10-03 AEST。当前精确模型及跨轮token未核实。没有真机抓包或部署；先完成修改与自动回归。

网关在单账号/号池选择后复用 fwd_identity 生成握手映射与内部 x-bf-codex-identity 上下文。客户端传入同名头先剥除，上下文只含原会话、映射会话及公司安装号，不含token。Bifrost的原生消息层利用已有 gjson/sjson 只修改指定字段，并剥除内部上下文头。无上下文时保持旧消息原样。客户端UA、模型、独立turn/request ID、业务输入和未知字段保留。

## 自动测试

测试使用隔离临时目录、本地假HTTP/WS上游、假凭据，无真实公司token或付费接口。

- `python -m pytest tests/test_bifrost_swap_proxy.py -k websocket -q`：初版1 passed；后续扩为单账号/号池两种路径，最终全量两项均通过。
- `python ops/test_all.py --md`：首轮381/382，已有错误回复头保留用例失败；单独复跑1 passed。第二轮383/383通过、0失败，用时125秒。首轮间歇失败保留记录。
- Bifrost工作树根目录先 `go work init ./core ./framework ./transports ./plugins/governance ./plugins/logging ./plugins/telemetry ./plugins/maxim ./plugins/otel ./plugins/modelcatalogresolver ./plugins/jsonparser ./plugins/routing ./plugins/semanticcache ./plugins/mocker ./plugins/compat ./plugins/prompts`，避免引用发布版模块；go.work为本机测试环境，不提交。
- `go test ./transports/bifrost-http/handlers ./transports/bifrost-http/integrations ./transports/bifrost-http/websocket -run 'TestCodexWSIdentity|TestChatGPT' -count=1`：handlers与integrations PASS，websocket该筛选无测试。
- 三包全量 `go test -json ... -count=1`：1308 pass / 104 fail，含子测试；与保存的历史未修改基线104个失败名称集合相同，无新增失败名称。基线版本不是本轮同提交重跑，不冒称严格同版基线。Windows/CGO/SQLite等旧环境失败保留。

新增回归验证真实本地WS多轮头/消息一致、内部头不泄漏、限额、断连、未知字节和大整数保留、业务输入及独立turn ID不改、映射稳定及账号隔离、非法上下文拒绝。

原始自动输出位于本机 aigw-development/ws-identity67-*；抓包原件仍在 captures，不进入Git。本项代码与测试推送后继续其他需求，最后统一真机抓包；整个FWD-67未完成。本项开发开始精确时间及累计token未核实；测试时段2026-10-03 23:26–23:32 AEST（报告时间至结束），不是整个任务累计时长。
