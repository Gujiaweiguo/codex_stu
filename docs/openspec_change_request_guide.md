# OpenSpec 变更需求（Change）帮助文档（Codex CLI）

这份文档用于指导你在 **Codex CLI** 环境下创建并推进一次 OpenSpec 变更（change）：把需求写清、让 Codex 小步实现、并用可复制的命令完成验证与归档。

> 适用场景：你希望把一次需求/改动“规格化”成可审阅的 proposal + tasks，并按步骤落地。

---

## 1) 前置条件与受限环境提示

### 前置条件
- Node.js（若你要使用 OpenSpec CLI）
- Codex CLI（用于对话驱动与执行 `/openspec-*` 命令）

### 受限环境（常见于 Codex CLI 容器）
如果你看到类似配置：
- `network_access=restricted`：安装 OpenSpec CLI（例如 `npm install -g ...`）可能需要你手动批准/允许联网。
- `approval_policy=on-request`：涉及联网/潜在破坏性命令时，Codex 会请求你确认。

建议做法：
- 任何“安装依赖/联网”先让 Codex **只给命令与影响说明**，你确认后再执行。
- 变更本身（proposal/tasks）可以先写在仓库里，落地实现再逐步执行验证。

---

## 2) 初始化 OpenSpec（第一次使用时）

### 2.0 建议先确认是否有 `AGENTS.md`（没有就先补一份）
`AGENTS.md` 不是 OpenSpec 的硬性要求，但它能让 Codex 明确“这个仓库怎么协作”，避免改动越界（例如：默认只改 Markdown、不要自动 commit/push、不要写敏感信息等）。

先检查仓库根目录是否已有：
```bash
test -f AGENTS.md && echo "AGENTS.md exists" || echo "AGENTS.md missing"
```

如果缺失，建议先创建一个最小版 `AGENTS.md`（示例，按你的仓库习惯改）：
```markdown
<INSTRUCTIONS>
# 项目协作约定（最小版）
- 默认只修改 docs/ 与图片资源，除非我明确要求改代码
- 不要引入新依赖/新工具链，除非我批准
- 不要自动 git commit / git push
- 示例不要写入密码、token、私钥等敏感信息（用占位符或环境变量）
- 输出命令必须可复制，执行前说明影响与回滚方式
</INSTRUCTIONS>
```

在仓库根目录执行（需要 OpenSpec CLI；若未安装请先安装）：
```bash
openspec init
```

初始化完成后通常会生成 `openspec/` 目录（以及可能的项目级说明文件）。你可以先确认目录是否出现：
```bash
ls -la openspec
```

---

## 3) 创建一个 change（变更需求）

### 命名建议
- 用 `kebab-case`：`fix-login-redirect`、`add-rbac-audit-log`
- 名字表达“要做什么”，避免过泛：`refactor`、`update` 这类不要用

### 创建方式 A：用 Codex 的 `/openspec-*` 命令（若你的环境支持）
在对话框里输入：
```text
/openspec-proposal <change-name>
```

### 创建方式 B：用 OpenSpec CLI（以你本机 `openspec --help` 输出为准）
先查看可用子命令：
```bash
openspec --help
```

创建 change 的命令可能因版本不同而略有差异；如果你不确定：
- 先用方式 A（`/openspec-proposal`）
- 或把 `openspec --help` 的输出贴给 Codex，让它给出最小可行命令

---

## 4) 写好 proposal：让“变更可执行、可验收”

创建 change 后，你会在 `openspec/` 下看到对应的变更文件（用 `openspec list` / `openspec show` 定位最可靠）。

常用检查命令：
```bash
openspec list
openspec show <change-name>
```

### proposal 推荐模板（可直接复制）
把下面内容填到你的 change proposal（具体文件路径以 `openspec show` 为准）：
```markdown
# 变更标题：<一句话说明要做什么>

## 背景
- 为什么要改：当前问题/限制/风险是什么？

## 目标
- [ ] 目标 1：
- [ ] 目标 2：

## 非目标（明确不做什么）
- 不做：

## 范围
- 影响的模块/目录：
- 明确不改的目录/文件：

## 方案概述
- 方案要点（最多 5 条）：

## 任务拆分（小步可交付）
1) 第 1 步：……（改哪些文件）
2) 第 2 步：……

## 验收标准
- 正常场景：
- 错误/边界场景：

## 测试计划（必须可复制）
```bash
# 例：只写你真正要跑的命令
<command-1>
<command-2>
```

## 风险与回滚
- 风险：
- 回滚方式（例如：回退哪些文件/禁用开关）：
```

写 proposal 的最低标准（不满足就先别开始实现）：
- 有清晰的“范围/非范围”
- 有可执行的“测试计划”（能在终端复制运行）
- tasks 拆到每一步都能独立验证/审阅

---

## 5) 让 Codex 按 tasks 小步落地（推荐话术）

### 先只要计划（不动手）
```text
先不要写代码。请阅读当前 change 的 proposal/tasks（用 openspec show 的内容为准），列出你需要澄清的 1–3 个问题。
然后给出 4–6 步执行计划：每一步写清要改哪些文件、怎么验证（终端命令）、预期结果。
```

### 开始执行第 N 步（强约束，防越界）
```text
开始第 N 步。这一步只允许改这些文件：<填文件列表>。
不要引入新依赖；不要移动/重命名文件；不要做无关重构。
改完只交付：改动摘要 + 我看 diff 的 3 个点 + 终端验证命令（含预期成功标准）。
```

---

## 6) 校验与归档

每完成一轮 tasks，建议先校验 change 结构与一致性：
```bash
openspec validate <change-name>
```

完成全部任务并验证通过后，再归档（会把变更合并回 specs/基线，具体行为以你的 OpenSpec 配置为准）：
```bash
openspec archive <change-name> --yes
```

---

## 7) 常见问题（排雷）

- change 太大：把 tasks 拆细，每一步必须能验证；否则会变成“不可审阅的大改”。
- 没写测试计划：先补最小验证命令，再开始实现。
- 联网/安装被卡：把 `openspec --help`、失败输出贴给 Codex，让它给“离线替代方案”或让你明确批准再做。
- 变更越界：每一步都限定文件列表；发现越界就要求回滚并重做。
