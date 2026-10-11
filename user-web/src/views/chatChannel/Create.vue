<template>
  <div class="chat-channel-form-page">
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <h2>{{ $t('auto.Create.25754973') }}</h2>
        <el-button @click="goBack">{{ $t('返回') }}</el-button>
      </div>
    </el-card>

    <el-card shadow="never">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="140px" style="max-width: 720px">
        <el-form-item ::label="$t('auto.Create.bbeea931')" prop="channel_name">
          <el-input v-model="form.channel_name" ::placeholder="$t('auto.Create.1fca53a4')" maxlength="100" show-word-limit />
        </el-form-item>

        <el-form-item ::label="$t('auto.Create.aa082b8f')" prop="allowed_origins">
          <el-input
            v-model="originsText"
            type="textarea"
            :rows="4"
            :placeholder="$t('auto.Create.92fc876b')"
          />
          <div class="form-tip">{{ $t('auto.Create.81e2beec') }}<code>*</code>{{ $t('auto.Create.4a66e4cc') }}</div>
        </el-form-item>

        <el-form-item ::label="$t('auto.Create.3a47ede0')">
          <el-input
            v-model="form.welcome_message"
            type="textarea"
            :rows="3"
            ::placeholder="$t('auto.Create.364e2934')"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>

        <el-form-item ::label="$t('auto.Create.648b4861')">
          <el-color-picker v-model="form.widget_color" />
          <span class="form-tip">{{ $t('聊天窗主色，建议与网站主色一致') }}</span>
        </el-form-item>

        <el-form-item ::label="$t('auto.Create.cd7a699c')">
          <el-radio-group v-model="form.widget_position">
            <el-radio value="bottom-right">{{ $t('右下角') }}</el-radio>
            <el-radio value="bottom-left">{{ $t('左下角') }}</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item ::label="$t('auto.Create.e9bfd4b7')">
          <el-input v-model="form.widget_title" ::placeholder="$t('auto.Create.6012f364')" maxlength="50" />
        </el-form-item>

        <el-form-item ::label="$t('auto.Create.da80092f')">
          <el-switch v-model="form.auto_assign" />
          <span class="form-tip">{{ $t('关闭后所有访客消息需人工接管') }}</span>
        </el-form-item>

        <el-form-item ::label="$t('auto.Create.f2297685')">
          <el-input-number
            v-model="form.confidence_threshold"
            :min="0"
            :max="1"
            :step="0.05"
            :precision="2"
          />
          <span class="form-tip">{{ $t('auto.Create.662f103d') }}</span>
        </el-form-item>

        <el-form-item ::label="$t('auto.Create.1ca3f6d2')">
          <el-select
            v-model="form.default_agent_id"
            :placeholder="$t('auto.Create.611251bd')"
            clearable
            filterable
            style="width: 100%"
            :loading="loadingAgents"
          >
            <el-option
              v-for="a in enabledAgents"
              :key="a.id"
              :label="`[${a.agent_code || '#' + a.id}] ${a.name}`"
              :value="a.id"
            >
              <div style="display: flex; justify-content: space-between; align-items: center">
                <span>{{ a.name }}</span>
                <el-tag size="small" type="info">{{ getAgentTypeLabel(a.agent_type) }}</el-tag>
              </div>
            </el-option>
          </el-select>
          <div class="form-tip">{{ $t('auto.Create.4159d328') }}</div>
        </el-form-item>

        <el-form-item ::label="$t('auto.Create.07ae792a')" prop="target_language">
          <el-select v-model="form.target_language" ::placeholder="$t('auto.Create.a28232ae')" clearable style="width: 100%">
            <el-option
              v-for="lang in targetLanguageOptions"
              :key="lang.value"
              :label="lang.label"
              :value="lang.value"
            />
          </el-select>
          <div class="form-tip">{{ $t('auto.Create.ac178471') }}</div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="saving" @click="onSubmit">{{ $t('创建并获取凭证') }}</el-button>
          <el-button @click="goBack">{{ $t('取消') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    
    <el-dialog
      v-model="credentialsVisible"
      ::title="$t('auto.Create.e049c740')"
      width="640px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
    >
      <el-alert type="warning" :closable="false" style="margin-bottom: 16px">
        <p style="margin: 0 0 4px"><strong>{{ $t('auto.Create.b5e59faf') }}</strong></p>
        <p style="margin: 0">{{ $t('auto.Create.67a7c570') }}</p>
      </el-alert>

      <el-form label-width="120px">
        <el-form-item ::label="$t('auto.Create.88145263')">
          <code class="cred">{{ createdData?.Channel?.channel_id }}</code>
        </el-form-item>
        <el-form-item label="AppKey">
          <code class="cred">{{ createdData?.AppKey }}</code>
          <el-button link type="primary" @click="copy(createdData?.AppKey)">
            <el-icon><CopyDocument /></el-icon>{{ $t('复制') }}
          </el-button>
        </el-form-item>
        <el-form-item label="AppSecret">
          <code class="cred secret">{{ createdData?.AppSecret }}</code>
          <el-button link type="primary" @click="copy(createdData?.AppSecret)">
            <el-icon><CopyDocument /></el-icon>{{ $t('复制') }}
          </el-button>
        </el-form-item>
      </el-form>

      <el-divider />
      <h4>{{ $t('集成代码（嵌入到企业网站）') }}</h4>
      <pre class="code-block">{{ embedCode }}</pre>
      <el-button type="primary" @click="copy(embedCode)">{{ $t('复制集成代码') }}</el-button>

      <template #footer>
        <el-button type="primary" @click="onDone">{{ $t('已完成保存') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { CopyDocument } from '@element-plus/icons-vue'
import { createChannel } from '@/api/chatChannel'
import { listEnabledAgents } from '@/api/aiAgent'
import { createBinding } from '@/api/channelAgentBinding'
import { TARGET_LANGUAGE_OPTIONS } from '@/constants/languages'

const router = useRouter()
const { t } = useI18n()

const formRef = ref()
const saving = ref(false)
const credentialsVisible = ref(false)
const createdData = ref(null)
const originsText = ref('')

const enabledAgents = ref([]);
const loadingAgents = ref(false)

const getAgentTypeLabel = (type) => {
  const map = { sales: t('auto.Create.7f0dc5e3'), customer_service: t('auto.Create.e7dea7c4'), hybrid: t('auto.Create.41ce131a') }
  return map[type] || type || '-'
}

const form = ref({
  channel_name: '',
  allowed_origins: [],
  welcome_message: t('auto.Create.0b7281ab'),
  widget_color: '#1989fa',
  widget_position: 'bottom-right',
  widget_title: t('auto.Create.7f72e8fe'),
  auto_assign: true,
  confidence_threshold: 0.70,
  default_agent_id: null
})

const targetLanguageOptions = TARGET_LANGUAGE_OPTIONS

const rules = {
  channel_name: [{ required: true, message: i18n.global.t(t('auto.Create.3fdb506c')), trigger: 'blur' }],
  allowed_origins: [
    {
      validator: (rule, value, cb) => {
        const lines = originsText.value.split('\n').map(s => s.trim()).filter(Boolean)
        if (lines.length === 0) {
          cb(new Error(t('auto.Create.c232ed81')))
        } else {
          cb()
        }
      },
      trigger: 'blur'
    }
  ]
}

const embedCode = computed(() => {
  if (!createdData.value?.app_key) return ''
  const baseURL = window.location.origin
  const closeTag = '<' + '/script>'
  return `<!-- 将以下代码嵌入到企业网站 </body> 之前 -->
<script src="${baseURL}/embed/marketing-chat-widget.iife.js" data-app-key="${createdData.value.app_key}">${closeTag}
<!-- 集成完成！刷新页面即可看到右下角浮标 -->`
})

const onSubmit = async () => {
  let valid
  try {
    valid = await formRef.value.validate()
  } catch (e) {
    valid = false
  }
  if (!valid) return
  saving.value = true
  try {
    form.value.allowed_origins = originsText.value
      .split('\n')
      .map(s => s.trim())
      .filter(Boolean)
    if (form.value.allowed_origins.includes('*')) {
      ElMessage.warning(i18n.global.t(t('auto.Create.7c583b4b')))
    }
    const payload = { ...form.value }
    delete payload.default_agent_id;
    const res = await createChannel(payload)
    createdData.value = res
    const channelId = res?.channel_id || res?.Channel?.channel_id;
    if (channelId && form.value.default_agent_id) {
      try {
        await createBinding({
          channel_type: 'web_embed',
          account_id: channelId,
          agent_id: form.value.default_agent_id,
          is_primary: true,
          enabled: true
        })
        ElMessage.success(i18n.global.t(t('auto.Create.7bcc8a4e')))
      } catch (bindErr) {
        ElMessage.warning(t('auto.Create.a2e0d831') + (bindErr?.message || t('auto.Create.974e7484')));
      }
    }
    credentialsVisible.value = true
  } catch (err) {
    ElMessage.error(t('auto.Create.3d6b61b7') + (err?.message || err))
  } finally {
    saving.value = false
  }
}

const onDone = () => {
  credentialsVisible.value = false
  router.push({ name: 'ChatChannelList' })
}

const goBack = () => router.push({ name: 'ChatChannelList' })

const copy = async (text) => {
  if (!text) return
  try {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(text)
    } else {
      throw new Error('clipboard unavailable')
    }
    ElMessage.success(i18n.global.t(t('auto.Create.52e6abbe')))
  } catch {
    try {
      const ta = document.createElement('textarea')
      ta.value = text
      ta.style.position = 'fixed'
      ta.style.opacity = '0'
      document.body.appendChild(ta)
      ta.select()
      document.execCommand('copy')
      document.body.removeChild(ta)
      ElMessage.success(i18n.global.t(t('auto.Create.52e6abbe')))
    } catch {
      ElMessage.error(i18n.global.t(t('auto.Create.063c173b')))
    }
  }
}

onMounted(() => {
  loadingAgents.value = true;
  listEnabledAgents()
    .then((res) => {
      const list = Array.isArray(res) ? res : res?.list || res?.items || []
      enabledAgents.value = list
    })
    .catch(() => {
      enabledAgents.value = []
    })
    .finally(() => {
      loadingAgents.value = false
    })
})
</script>

<style scoped>
.chat-channel-form-page {
  padding: 0;
}
.header-card {
  margin-bottom: 16px;
}
.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-content h2 {
  margin: 0;
  font-size: 20px;
}
.form-tip {
  margin-left: 12px;
  font-size: 12px;
  color: #909399;
}
.cred {
  font-family: 'SF Mono', Monaco, Consolas, monospace;
  font-size: 13px;
  background: #f5f7fa;
  padding: 4px 8px;
  border-radius: 3px;
  user-select: all;
  word-break: break-all;
}
.cred.secret {
  background: #fef0f0;
  color: #EF4444;
}
.code-block {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 12px 16px;
  border-radius: 4px;
  font-size: 12px;
  font-family: 'SF Mono', Monaco, Consolas, monospace;
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
