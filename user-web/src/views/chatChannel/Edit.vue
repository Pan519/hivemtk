<template>
  <div class="chat-channel-form-page">
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <h2>{{ $t('编辑客服渠道') }}</h2>
        <el-button @click="goBack">{{ $t('返回') }}</el-button>
      </div>
    </el-card>

    <el-card v-loading="loading" shadow="never">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="140px" style="max-width: 720px">
        <el-form-item ::label="$t('auto.Edit.88145263')">
          <code class="readonly">{{ form.channel_id }}</code>
        </el-form-item>
        <el-form-item label="AppKey">
          <code class="readonly">{{ form.app_key }}</code>
          <span class="form-tip">{{ $t('auto.Edit.b3fb7972') }}</span>
        </el-form-item>
        <el-form-item ::label="$t('auto.Edit.bbeea931')" prop="channel_name">
          <el-input v-model="form.channel_name" maxlength="100" show-word-limit />
        </el-form-item>
        <el-form-item ::label="$t('auto.Edit.aa082b8f')" prop="allowed_origins_text">
          <el-input
            v-model="originsText"
            type="textarea"
            :rows="4"
            ::placeholder="$t('auto.Edit.33f6e100')"
          />
        </el-form-item>
        <el-form-item ::label="$t('auto.Edit.1b857887')">
          <el-radio-group v-model="form.status">
            <el-radio value="active">{{ $t('启用') }}</el-radio>
            <el-radio value="disabled">{{ $t('禁用') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item ::label="$t('auto.Edit.3a47ede0')">
          <el-input v-model="form.welcome_message" type="textarea" :rows="3" maxlength="500" show-word-limit />
        </el-form-item>
        <el-form-item ::label="$t('auto.Edit.648b4861')">
          <el-color-picker v-model="form.widget_color" />
        </el-form-item>
        <el-form-item ::label="$t('auto.Edit.cd7a699c')">
          <el-radio-group v-model="form.widget_position">
            <el-radio value="bottom-right">{{ $t('右下角') }}</el-radio>
            <el-radio value="bottom-left">{{ $t('左下角') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item ::label="$t('auto.Edit.e9bfd4b7')">
          <el-input v-model="form.widget_title" maxlength="50" />
        </el-form-item>
        <el-form-item ::label="$t('auto.Edit.da80092f')">
          <el-switch v-model="form.auto_assign" />
        </el-form-item>
        <el-form-item ::label="$t('auto.Edit.f2297685')">
          <el-input-number v-model="form.confidence_threshold" :min="0" :max="1" :step="0.05" :precision="2" />
        </el-form-item>
        <el-form-item ::label="$t('auto.Edit.1ca3f6d2')">
          <el-select
            v-model="form.default_agent_id"
            :placeholder="$t('auto.Edit.43205678')"
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
          <div class="form-tip">{{ $t('auto.Edit.09b6f046') }}</div>
        </el-form-item>
        <el-form-item ::label="$t('auto.Edit.07ae792a')" prop="target_language">
          <el-select v-model="form.target_language" ::placeholder="$t('auto.Edit.a28232ae')" clearable style="width: 100%">
            <el-option
              v-for="lang in targetLanguageOptions"
              :key="lang.value"
              :label="lang.label"
              :value="lang.value"
            />
          </el-select>
          <div class="form-tip">{{ $t('auto.Edit.ac178471') }}</div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving" @click="onSubmit">{{ $t('保存') }}</el-button>
          <el-button @click="goBack">{{ $t('取消') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getChannel, updateChannel } from '@/api/chatChannel'
import { listEnabledAgents } from '@/api/aiAgent'
import { listBindings, createBinding, updateBinding, deleteBinding } from '@/api/channelAgentBinding'
import { TARGET_LANGUAGE_OPTIONS } from '@/constants/languages'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const formRef = ref()
const loading = ref(false)
const saving = ref(false)
const originsText = ref('')
const targetLanguageOptions = TARGET_LANGUAGE_OPTIONS

const enabledAgents = ref([]);
const loadingAgents = ref(false)
const initialDefaultAgentId = ref(null);

const getAgentTypeLabel = (type) => {
  const map = { sales: t('auto.Edit.7f0dc5e3'), customer_service: t('auto.Edit.e7dea7c4'), hybrid: t('auto.Edit.41ce131a') }
  return map[type] || type || '-'
}

const form = ref({
  channel_id: '',
  app_key: '',
  channel_name: '',
  status: 'active',
  welcome_message: '',
  widget_color: '#1989fa',
  widget_position: 'bottom-right',
  widget_title: t('auto.Edit.7f72e8fe'),
  auto_assign: true,
  confidence_threshold: 0.70,
  target_language: '',
  default_agent_id: null
})

const rules = {
  channel_name: [{ required: true, message: i18n.global.t(t('auto.Edit.3fdb506c')), trigger: 'blur' }]
}

const loadDetail = async () => {
  loading.value = true
  try {
    const res = await getChannel(route.params.id)
    Object.assign(form.value, res)
    form.value.target_language = res?.target_language || '';
    originsText.value = (res.allowed_origins || '').split(',').filter(Boolean).join('\n')
    await Promise.all([loadEnabledAgents(), loadCurrentDefaultAgent()]);
  } catch (err) {
    ElMessage.error(t('auto.Edit.c4fa9b74') + (err?.message || err))
  } finally {
    loading.value = false
  }
}

const loadEnabledAgents = async () => {
  loadingAgents.value = true
  try {
    const res = await listEnabledAgents()
    const list = Array.isArray(res) ? res : res?.list || res?.items || []
    enabledAgents.value = list
  } catch {
    enabledAgents.value = []
  } finally {
    loadingAgents.value = false
  }
}

const loadCurrentDefaultAgent = async () => {
  try {
    const res = await listBindings({ channel_type: 'web_embed', account_id: form.value.channel_id });
    const list = Array.isArray(res) ? res : res?.list || res?.items || []
    const def = list.find((b) => b.is_default || b.is_primary) || list[0]
    if (def) {
      const aid = def.agent_id || def.AgentId
      form.value.default_agent_id = aid
      initialDefaultAgentId.value = aid
    }
  } catch { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}

const syncDefaultAgentBinding = async () => {
  const cur = form.value.default_agent_id ? Number(form.value.default_agent_id) : null
  const orig = initialDefaultAgentId.value ? Number(initialDefaultAgentId.value) : null
  if (cur === orig) return
  if (orig) {
    try {
      const res = await listBindings({ channel_type: 'web_embed', account_id: form.value.channel_id, agent_id: orig })
      const list = Array.isArray(res) ? res : res?.list || res?.items || []
      for (const b of list) {
        if (b.is_default || b.is_primary) {
          await deleteBinding(b.id || b.ID).catch(() => null)
        }
      }
    } catch { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
  }
  if (cur) {
    await createBinding({
      channel_type: 'web_embed',
      account_id: form.value.channel_id,
      agent_id: cur,
      is_primary: true,
      enabled: true
    })
    initialDefaultAgentId.value = cur
  } else {
    initialDefaultAgentId.value = null
  }
}

const onSubmit = async () => {
  try {
    await formRef.value.validate()
  } catch {
    return;
  }
  saving.value = true
  try {
    const payload = {
      channel_name: form.value.channel_name,
      status: form.value.status,
      allowed_origins: originsText.value.split('\n').map(s => s.trim()).filter(Boolean),
      welcome_message: form.value.welcome_message,
      widget_color: form.value.widget_color,
      widget_position: form.value.widget_position,
      widget_title: form.value.widget_title,
      auto_assign: form.value.auto_assign,
      confidence_threshold: form.value.confidence_threshold,
      target_language: form.value.target_language || ''
    }
    await updateChannel(form.value.channel_id, payload)
    try {
      await syncDefaultAgentBinding()
    } catch (bindErr) {
      ElMessage.warning(t('auto.Edit.c2669cb5') + (bindErr?.message || t('auto.Edit.974e7484')))
    }
    ElMessage.success(i18n.global.t(t('auto.Edit.f8dfedcd')))
    router.push({ name: 'ChatChannelList' })
  } catch (err) {
    ElMessage.error(t('auto.Edit.78d235cc') + (err?.message || err))
  } finally {
    saving.value = false
  }
}

const goBack = () => router.push({ name: 'ChatChannelList' })

onMounted(loadDetail)
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
.readonly {
  font-family: 'SF Mono', Monaco, Consolas, monospace;
  font-size: 13px;
  background: #f5f7fa;
  padding: 4px 8px;
  border-radius: 3px;
  user-select: all;
}
.form-tip {
  margin-left: 12px;
  font-size: 12px;
  color: #909399;
}
</style>
