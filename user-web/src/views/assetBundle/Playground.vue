<template>
  <div class="playground">
    
    <el-card class="header-card" shadow="never">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-form-item :label="$t('auto.Playground.69847d0a')" label-width="80px">
            <el-input v-model="bundle.asset_id" placeholder="hive_sales_vape_cn_001" :disabled="isEdit" />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item :label="$t('auto.Playground.0acf6dc2')" label-width="90px">
            <el-input v-model="bundle.title" :placeholder="$t('auto.Playground.8774f184')" />
          </el-form-item>
        </el-col>
        <el-col :span="4">
          <el-form-item :label="$t('auto.Playground.fe2df04a')" label-width="60px">
            <el-input v-model="bundle.version" placeholder="1.0.0" />
          </el-form-item>
        </el-col>
        <el-col :span="4">
          <el-form-item :label="$t('auto.Playground.6adafe27')" label-width="60px">
            <el-input v-model="bundle.industry" :placeholder="$t('auto.Playground.4529dfda')" />
          </el-form-item>
        </el-col>
        <el-col :span="4">
          <div class="header-actions">
            <el-button type="primary" :loading="saving" @click="handleSave">{{ $t('auto.Playground.afcc3e34') }}</el-button>
            <el-button v-if="bundle.id" type="success" @click="handlePublish">{{ $t('auto.Playground.e822c24e') }}</el-button>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <el-row :gutter="16" class="main-row">
      
      <el-col :span="14">
        <el-card class="left-panel" shadow="never">
          <template #header>
            <div class="panel-header">
              <span>{{ $t('auto.Playground.709de71c') }}</span>
              <el-button type="primary" size="small" @click="addMessage('user')">{{ $t('auto.Playground.581c3818') }}</el-button>
            </div>
          </template>

          <div class="messages-list">
            <div
              v-for="(msg, idx) in bundle.messages"
              :key="idx"
              class="message-card"
              :class="'role-' + msg.role"
            >
              <div class="message-header">
                <el-tag :type="roleTagType(msg.role)" size="small">
                  {{ roleLabel(msg.role) }}
                </el-tag>
                <span v-if="idx === 0 && msg.role === 'system'" class="hint">{{ $t('auto.Playground.1bac1244') }}</span>
                <span v-else-if="msg.role === 'user'" class="hint">{{ $t('auto.Playground.d5eb2ea7') }}</span>
                <span v-else-if="msg.role === 'assistant'" class="hint">{{ $t('auto.Playground.26ff4a38') }}</span>
                <div class="message-actions">
                  <el-button
                    v-if="idx > 0"
                    size="small"
                    link
                    @click="moveUp(idx)"
                  >↑</el-button>
                  <el-button
                    v-if="idx < bundle.messages.length - 1"
                    size="small"
                    link
                    @click="moveDown(idx)"
                  >↓</el-button>
                  <el-button
                    v-if="idx > 0"
                    size="small"
                    type="danger"
                    link
                    @click="removeMessage(idx)"
                  >{{ $t('auto.Playground.2f4aaddd') }}</el-button>
                </div>
              </div>
              <el-input
                v-model="msg.content"
                type="textarea"
                :rows="idx === 0 ? 10 : 4"
                :placeholder="rolePlaceholder(msg.role)"
                class="message-textarea"
              />
            </div>
          </div>

          <div class="add-pair-row">
            <el-button @click="addMessage('user')">{{ $t('auto.Playground.e043f744') }}</el-button>
            <el-button @click="addMessage('assistant')">{{ $t('auto.Playground.9ce2e580') }}</el-button>
            <el-button @click="addMessage('system')">{{ $t('auto.Playground.cf57fdb0') }}</el-button>
          </div>
        </el-card>
      </el-col>

      
      <el-col :span="10">
        <el-card class="right-panel" shadow="never">
          <template #header>
            <div class="panel-header">
              <span>{{ $t('auto.Playground.02ba7a39') }}</span>
            </div>
          </template>

          
          <div class="config-section">
            <div class="section-title">{{ $t('auto.Playground.d20472c2') }}</div>
            <el-form label-width="120px" size="small">
              <el-form-item :label="$t('auto.Playground.610269e4')">
                <el-input v-model="sandbox.endpoint" placeholder="http://localhost:11434/v1/chat/completions" />
              </el-form-item>
              <el-form-item :label="$t('auto.Playground.0e0f4df3')">
                <el-input v-model="sandbox.model" placeholder="qwen2.5:7b" />
              </el-form-item>
              <el-form-item :label="$t('auto.Playground.5b8c65d8')">
                <el-slider v-model="sandbox.temperature" :min="0" :max="1" :step="0.1" show-input />
              </el-form-item>
              <el-form-item :label="$t('auto.Playground.e09c9b51')">
                <el-input-number v-model="sandbox.maxTokens" :min="128" :max="8192" :step="128" />
              </el-form-item>
            </el-form>
          </div>

          
          <div class="chat-section">
            <div class="section-title">{{ $t('auto.Playground.92bdbbc2') }}</div>
            <div class="chat-history" ref="chatHistoryRef">
              <div
                v-for="(m, i) in sandbox.history"
                :key="i"
                class="chat-msg"
                :class="'chat-' + m.role"
              >
                <div class="chat-role">{{ roleLabel(m.role) }}</div>
                <pre class="chat-content">{{ m.content }}</pre>
              </div>
              <el-empty v-if="!sandbox.history.length" :description="$t('auto.Playground.53a0637a')" :image-size="60" />
            </div>
            <div class="chat-input-row">
              <el-input
                v-model="sandbox.input"
                type="textarea"
                :rows="2"
                :placeholder="$t('auto.Playground.1ba6a52e')"
                @keydown.enter.ctrl="runSandbox"
              />
              <el-button type="primary" :loading="sandbox.running" @click="runSandbox">
                ▶️ 运行
              </el-button>
            </div>
          </div>

          
          <div class="intent-section" v-if="sandbox.lastIntent">
            <div class="section-title">{{ $t('auto.Playground.6d5f3689') }}</div>
            <div class="intent-status">
              <el-tag :type="sandbox.lastIntent.valid ? 'success' : 'danger'">
                {{ sandbox.lastIntent.valid ? '🟢 合法格式，成功抓取' : '🔴 非法格式' }}
              </el-tag>
              <span v-if="sandbox.lastIntent.intent" class="intent-label">
                intent: <code>{{ sandbox.lastIntent.intent }}</code>
              </span>
              <span v-if="sandbox.lastIntent.durationMs != null" class="intent-label">
                耗时: {{ sandbox.lastIntent.durationMs }} ms
              </span>
            </div>
            <pre class="intent-json">{{ sandbox.lastIntent.rawJSON || '(无 JSON 块)' }}</pre>
            <div v-if="sandbox.lastIntent.capturedData && Object.keys(sandbox.lastIntent.capturedData).length" class="captured-data">
              <div class="section-subtitle">captured_data:</div>
              <pre>{{ JSON.stringify(sandbox.lastIntent.capturedData, null, 2) }}</pre>
            </div>
          </div>

          
          <div class="woven-section" v-if="sandbox.wovenMessages && sandbox.wovenMessages.length">
            <div class="section-title">🧶 织布后的最终 messages 数组（{{ sandbox.wovenMessages.length }} 条）</div>
            <pre class="woven-preview">{{ JSON.stringify(sandbox.wovenMessages, null, 2) }}</pre>
          </div>

          
          <div class="publish-section">
            <div class="section-title">{{ platformEnabled ? '💰 生态上架配置' : '📦 本地发布' }}</div>
            <el-form
              v-if="platformEnabled"
              label-width="120px"
              size="small"
            >
              <el-form-item :label="$t('auto.Playground.694a76ec')">
                <el-input v-model="bundle.price" placeholder="$ 299.00" />
              </el-form-item>
              <el-form-item :label="$t('auto.Playground.4705b884')">
                <el-radio-group v-model="bundle.scope">
                  <el-radio label="private">{{ $t('auto.Playground.3dc5185d') }}</el-radio>
                  <el-radio label="shared">{{ $t('auto.Playground.fffd3ce7') }}</el-radio>
                  <el-radio label="official">{{ $t('auto.Playground.f5c3795b') }}</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-form>
            <el-button type="primary" plain class="publish-btn" @click="handlePublish">
              {{ platformEnabled ? '🚀 审核上架到官方蜂巢商城' : '🚀 发布到本地资产库' }}
              <span class="publish-hint">{{ $t('auto.Playground.cff3361b') }}</span>
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createBundle, updateBundle, getBundleByAssetID,
  weaveBundle, publishBundle, submitToPlatform
} from '@/api/assetBundle'
import { useUserStore } from '@/stores/user'
import { SystemApi } from '@/api/system'

const userStore = useUserStore()

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const chatHistoryRef = ref(null)
const saving = ref(false)
// 关态下后端对上架写接口回 403，入口直接不渲染；取不到就按关态处理
const platformEnabled = ref(false)

const aid = computed(() => route.params.aid || '')
const isEdit = computed(() => !!aid.value)

const bundle = reactive({
  id: null,
  asset_id: '',
  title: '',
  version: '1.0.0',
  industry: '',
  author: '',
  scope: 'private',
  status: 'draft',
  price: '',
  messages: [
    {
      role: 'system',
      content: '# 核心角色与销冠人设\n你是一名经过严格训练、结果导向的【王牌私域销售代表】。\n\n# 强制业务结算协议\n你必须在每一次回复的【最后】附带一个 ```json 块：\n```json\n{\n  "intent": "faq|lead_capture|human_transfer",\n  "captured_data": {}\n}\n```'
    }
  ]
});

const sandbox = reactive({
  endpoint: 'http://localhost:11434/v1/chat/completions',
  model: 'qwen2.5:7b',
  temperature: 0.2,
  maxTokens: 1024,
  input: '',
  history: [],
  running: false,
  lastIntent: null,
  wovenMessages: null
});

const roleLabel = (r) => ({
  system: 'system',
  user: 'user',
  assistant: 'assistant',
  tool: 'tool'
}[r] || r)

const roleTagType = (r) => ({
  system: 'danger',
  user: 'primary',
  assistant: 'success',
  tool: 'warning'
}[r] || 'info')

const rolePlaceholder = (r) => ({
  system: t('auto.Playground.d282e812'),
  user: t('auto.Playground.97423adf'),
  assistant: 'AI 标准回复（含末尾 ```json 业务结算块```）',
  tool: t('auto.Playground.9742cd7b')
}[r] || '')

const addMessage = (role = 'user') => {
  bundle.messages.push({ role, content: '' })
  if (role === 'user') {
    bundle.messages.push({ role: 'assistant', content: '' });
  }
}

const removeMessage = (idx) => {
  bundle.messages.splice(idx, 1)
}

const moveUp = (idx) => {
  if (idx <= 0) return
  const arr = bundle.messages
  ;[arr[idx - 1], arr[idx]] = [arr[idx], arr[idx - 1]]
}

const moveDown = (idx) => {
  const arr = bundle.messages
  if (idx >= arr.length - 1) return
  ;[arr[idx], arr[idx + 1]] = [arr[idx + 1], arr[idx]]
}

const loadBundle = async () => {
  if (!aid.value) return
  try {
    const resp = await getBundleByAssetID(aid.value)
    const data = resp?.data || resp
    if (data) {
      bundle.id = data.id || null
      bundle.asset_id = data.asset_id || ''
      bundle.title = data.title || ''
      bundle.version = data.version || '1.0.0'
      bundle.industry = data.industry || ''
      bundle.author = data.author || ''
      bundle.scope = data.scope || 'private'
      bundle.status = data.status || 'draft'
      bundle.messages = (data.messages || []).map(m => ({ ...m }))
    }
  } catch (e) {
    ElMessage.error(t('auto.Playground.2cfc148d') + (e?.message || e))
  }
};

const handleSave = async () => {
  if (!bundle.asset_id) {
    ElMessage.warning(t('auto.Playground.b9b9a97e'))
    return
  }
  if (!bundle.title) {
    ElMessage.warning(t('auto.Playground.ac8238fb'))
    return
  }
  if (!bundle.messages.length) {
    ElMessage.warning(t('auto.Playground.a397f8ff'))
    return
  }
  saving.value = true
  try {
    const payload = {
      asset_id: bundle.asset_id,
      title: bundle.title,
      version: bundle.version,
      industry: bundle.industry,
      author: bundle.author || userStore.username,
      scope: bundle.scope,
      messages: bundle.messages
    }
    if (bundle.id) {
      payload.id = bundle.id
      await updateBundle(bundle.id, payload)
      ElMessage.success(t('auto.Playground.55aa6366'))
    } else {
      const resp = await createBundle(payload)
      const data = resp?.data || resp
      if (data && data.id) {
        bundle.id = data.id
        ElMessage.success(t('auto.Playground.3b108349'))
        router.replace(`/asset-bundle/playground/${bundle.asset_id}`)
      } else {
        ElMessage.success(t('auto.Playground.3b108349'))
      }
    }
  } catch (e) {
    ElMessage.error(t('auto.Playground.40f90217') + (e?.message || e))
  } finally {
    saving.value = false
  }
};

const handlePublish = async () => {
  if (!bundle.id) {
    ElMessage.warning(t('auto.Playground.9558c46d'))
    return
  }
  const toPlatform = platformEnabled.value
  try {
    await ElMessageBox.confirm(
      toPlatform ? t('auto.Playground.40ec618c') : t('auto.Playground.325fd5f5'),
      t('auto.Playground.e83a256e'), { type: 'warning' })
    await publishBundle(bundle.id)
    bundle.status = 'active'
    if (!toPlatform) {
      ElMessage.success(t('auto.Playground.ac3aba2b'))
      return
    }
    try {
      await submitToPlatform(bundle.id)
      ElMessage.success(t('auto.Playground.1c603424'))
    } catch (e) {
      ElMessage.warning(t('auto.Playground.f716d0d8') + (e?.message || e))
    }
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error(t('auto.Playground.bd774adc') + (e?.message || e))
    }
  }
};

const JSON_BLOCK_RE = /```(?:json|JSON)?\s*(\{[\s\S]*?\})\s*```\s*$/;
const BARE_JSON_RE = /(\{"intent"\s*:[^{}]*(?:\{[^{}]*\}[^{}]*)*\})\s*$/

const extractIntent = (reply) => {
  const result = { valid: false, intent: '', capturedData: {}, rawJSON: '' }
  const blockRE = /```(?:json|JSON)?\s*(\{[\s\S]*?\})\s*```\s*$/;
  let m = reply.match(blockRE)
  if (!m) {
    const bareRE = /(\{"intent"\s*:[\s\S]*?\})\s*$/;
    m = reply.match(bareRE)
  }
  if (!m) return { ...result, reply }
  const rawJSON = m[1]
  result.rawJSON = rawJSON
  try {
    const parsed = JSON.parse(rawJSON)
    result.valid = true
    result.intent = parsed.intent || ''
    result.capturedData = parsed.captured_data || {}
    const cleaned = reply.slice(0, reply.indexOf(m[0])).replace(/\n+\s*$/, '');
    return { ...result, reply: cleaned }
  } catch (e) {
    return { ...result, reply }
  }
};

const runSandbox = async () => {
  if (!sandbox.input.trim()) {
    ElMessage.warning(t('auto.Playground.32f04b9e'))
    return
  }
  if (!bundle.messages.length) {
    ElMessage.warning(t('auto.Playground.0a3a6c8a'))
    return
  }
  sandbox.running = true
  sandbox.lastIntent = null
  sandbox.wovenMessages = null
  const userQuery = sandbox.input
  sandbox.history.push({ role: 'user', content: userQuery })
  sandbox.input = ''

  let weaveOk = false;
  try {
    const weaveResp = await weaveBundle({
      asset_id: bundle.asset_id,
      user_query: userQuery,
      sandbox: true,
      chat_history: sandbox.history.filter(m => m.role !== 'system').map(m => ({ role: m.role, content: m.content })),
      options: {
        rag_position: 'after_fewshots',
        max_history_messages: 10,
        strip_few_shot_json: false,
        include_merchant_vars: false
      }
    })
    weaveOk = true
    const weaveData = weaveResp?.data || weaveResp
    const wovenMessages = weaveData?.messages || []
    sandbox.wovenMessages = wovenMessages

    const llmResp = await fetch(sandbox.endpoint, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        model: sandbox.model,
        messages: wovenMessages,
        temperature: sandbox.temperature,
        max_tokens: sandbox.maxTokens,
        stream: false
      })
    });
    if (!llmResp.ok) {
      const errText = await llmResp.text()
      throw new Error(`LLM 返回 ${llmResp.status}: ${errText}`)
    }
    const llmData = await llmResp.json()
    const reply = llmData?.choices?.[0]?.message?.content || ''
    const startTs = Date.now()
    const extracted = extractIntent(reply);
    sandbox.history.push({ role: 'assistant', content: extracted.reply })
    sandbox.lastIntent = {
      valid: extracted.valid,
      intent: extracted.intent,
      capturedData: extracted.capturedData,
      rawJSON: extracted.rawJSON,
      durationMs: Date.now() - startTs
    }
  } catch (e) {
    if (!weaveOk) {
      const msg = (e && e.message) ? e.message : t('auto.Playground.80357119');
      ElMessage.error(t('auto.Playground.e32926c4') + msg)
      sandbox.history.push({
        role: 'assistant',
        content: `[织布失败] ${msg}\n\n提示：\n1. 请先「保存」资产包（织布需按 asset_id 读取已保存的 messages）\n2. 确认 asset_id 正确且未被删除`
      })
    } else {
      sandbox.history.push({
        role: 'assistant',
        content: `[LLM 调用失败] ${e?.message || e}\n\n提示：\n1. 请确认本地已启动 Ollama / vLLM\n2. 推理端地址正确（默认 http://localhost:11434/v1/chat/completions）\n3. CORS 已允许（Ollama 设置 OLLAMA_ORIGINS=*）`
      })
      ElMessage.error('LLM 调用失败')
    }
  } finally {
    sandbox.running = false
    await nextTick()
    if (chatHistoryRef.value) {
      chatHistoryRef.value.scrollTop = chatHistoryRef.value.scrollHeight
    }
  }
};

onMounted(() => {
  SystemApi.getInfo({ _silent: true })
    .then((info) => { platformEnabled.value = info?.platform_enabled === true })
    .catch(() => { platformEnabled.value = false })
  if (aid.value) {
    bundle.asset_id = aid.value
    loadBundle()
  }
})
</script>

<style scoped>
.playground {
  padding: 12px;
}
.header-card {
  margin-bottom: 12px;
}
.header-actions {
  display: flex;
  gap: 8px;
  padding-top: 4px;
}
.main-row {
  height: calc(100vh - 200px);
}
.left-panel, .right-panel {
  height: 100%;
  overflow: auto;
}
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}
.messages-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.message-card {
  border-left: 4px solid #409eff;
  padding: 8px 12px;
  background: #fafafa;
  border-radius: 4px;
}
.role-system {
  border-left-color: #f56c6c;
}
.role-user {
  border-left-color: #409eff;
}
.role-assistant {
  border-left-color: #67c23a;
}
.message-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}
.hint {
  font-size: 12px;
  color: #909399;
}
.message-actions {
  margin-left: auto;
}
.message-textarea :deep(.el-textarea__inner) {
  font-family: 'Courier New', monospace;
  font-size: 13px;
}
.add-pair-row {
  margin-top: 12px;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
.config-section, .chat-section, .intent-section, .woven-section, .publish-section {
  margin-bottom: 16px;
  padding: 12px;
  background: #fafafa;
  border-radius: 4px;
}
.section-title {
  font-weight: 600;
  margin-bottom: 10px;
  padding-bottom: 6px;
  border-bottom: 1px dashed #dcdfe6;
}
.section-subtitle {
  font-size: 12px;
  color: #606266;
  margin-top: 8px;
  margin-bottom: 4px;
}
.chat-history {
  max-height: 240px;
  overflow-y: auto;
  background: #fff;
  padding: 8px;
  border-radius: 4px;
  border: 1px solid #ebeef5;
  margin-bottom: 8px;
}
.chat-msg {
  margin-bottom: 8px;
  padding: 6px 8px;
  border-radius: 4px;
}
.chat-user {
  background: #ecf5ff;
  border-left: 3px solid #409eff;
}
.chat-assistant {
  background: #f0f9eb;
  border-left: 3px solid #67c23a;
}
.chat-system {
  background: #fef0f0;
  border-left: 3px solid #f56c6c;
}
.chat-role {
  font-size: 11px;
  color: #909399;
  margin-bottom: 2px;
}
.chat-content {
  font-size: 13px;
  white-space: pre-wrap;
  word-break: break-word;
  margin: 0;
  font-family: inherit;
}
.chat-input-row {
  display: flex;
  gap: 8px;
  align-items: flex-end;
}
.intent-status {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 6px;
}
.intent-label {
  font-size: 13px;
}
.intent-label code {
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: monospace;
}
.intent-json {
  background: #2d2d2d;
  color: #f8f8f2;
  padding: 8px;
  border-radius: 4px;
  font-size: 12px;
  font-family: 'Courier New', monospace;
  max-height: 160px;
  overflow: auto;
  margin: 0;
}
.captured-data pre {
  background: #f5f7fa;
  padding: 6px;
  border-radius: 4px;
  font-size: 12px;
  margin: 4px 0 0;
}
.woven-preview {
  background: #f5f7fa;
  padding: 8px;
  border-radius: 4px;
  font-size: 11px;
  max-height: 200px;
  overflow: auto;
  margin: 0;
}
.publish-btn {
  width: 100%;
}
.publish-hint {
  font-size: 11px;
  color: #909399;
  margin-left: 4px;
}
</style>
