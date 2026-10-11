<template>
  <div class="connectors-page">
    <PageHeader :title="$t('auto.Connectors.7de30cf7')" :subtitle="$t('auto.Connectors.e570eba0')" />

    <el-row :gutter="16">
      <el-col :span="6" v-for="c in connectors" :key="c.source">
        <el-card shadow="hover" class="connector-card">
          <div class="card-head">
            <span class="conn-name">{{ c.name }}</span>
            <el-tag :type="c.enabled ? 'success' : 'info'" size="small">{{ c.enabled ? '已启用' : '未启用' }}</el-tag>
          </div>
          <el-descriptions :column="1" size="small" class="conn-desc">
            <el-descriptions-item :label="$t('auto.Connectors.5f811dd8')">{{ c.configured ? '已配置' : '未配置' }}</el-descriptions-item>
            <el-descriptions-item v-if="c.last_test_at" :label="$t('auto.Connectors.a5990382')">
              <el-tag :type="c.last_test_ok ? 'success' : 'danger'" size="small">{{ c.last_test_ok ? '通过' : '失败' }}</el-tag>
              <span class="test-msg">{{ c.last_test_msg }}</span>
            </el-descriptions-item>
          </el-descriptions>

          <el-collapse accordion class="cred-form">
            <el-collapse-item :title="$t('auto.Connectors.d1df5468')">
              <el-form label-position="top" size="small">
                <el-form-item v-for="field in credFields(c.source)" :key="field.key" :label="field.label">
                  <el-input v-model="credForm[c.source][field.key]"
                    ::placeholder="$t('auto.Connectors.d96d8c87')"
                    show-password />
                </el-form-item>
                <el-switch v-model="credForm[c.source].enabled" :active-text="$t('auto.Connectors.7854b52a')" />
                <div class="form-actions">
                  <el-button size="small" type="primary" @click="saveCred(c.source)">{{ $t('auto.Connectors.bab96443') }}</el-button>
                  <el-button size="small" @click="testConn(c.source)">{{ $t('auto.Connectors.69e74756') }}</el-button>
                </div>
              </el-form>
            </el-collapse-item>
          </el-collapse>

          <div class="pull-box">
            <el-select v-model="pullKb[c.source]" size="small" :placeholder="$t('auto.Connectors.210f363d')" class="kb-select">
              <el-option v-for="kb in kbList" :key="kb.id" :label="kb.name" :value="String(kb.id)" />
            </el-select>
            <el-button size="small" type="success" :loading="pulling === c.source" :disabled="!c.enabled"
              @click="pull(c.source)">{{ $t('auto.Connectors.a2b59c0b') }}</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card v-if="lastResult" class="result-card">
      <template #header>最近拉取结果（{{ lastResult.source }}）：成功 {{ lastResult.imported }} / 失败 {{ lastResult.failed }} / 跳过 {{ lastResult.skipped }}</template>
      <el-table :data="lastResult.details" size="small" max-height="320">
        <el-table-column prop="title" :label="$t('auto.Connectors.59ceff46')" min-width="180" show-overflow-tooltip />
        <el-table-column prop="status" :label="$t('auto.Connectors.3fea7ca7')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'imported' ? 'success' : row.status === 'skipped' ? 'info' : 'danger'" size="small">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="message" :label="$t('auto.Connectors.f411d0f1')" min-width="200" show-overflow-tooltip />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/PageHeader.vue'
import { http } from '@/utils/request'

const { t } = useI18n()

const connectors = ref([])
const kbList = ref([])
const pulling = ref('')
const lastResult = ref(null)

const emptyCred = () => ({ enabled: false, token: '', app_id: '', app_secret: '', app_key: '', webhook_url: '', api_key: '' })
const credForm = ref({ notion: emptyCred(), feishu: emptyCred(), dingtalk: emptyCred(), crm: emptyCred() })
const pullKb = ref({})

const credFields = (source) => ({
  notion: [{ key: 'token', label: 'Internal Integration Token', placeholder: 'secret_xxx / ntn_xxx' }],
  feishu: [
    { key: 'app_id', label: 'App ID', placeholder: 'cli_xxx' },
    { key: 'app_secret', label: 'App Secret', placeholder: t('auto.Connectors.a86754cd') }
  ],
  dingtalk: [
    { key: 'app_key', label: 'AppKey', placeholder: t('auto.Connectors.7d9032dc') },
    { key: 'app_secret', label: 'AppSecret', placeholder: t('auto.Connectors.71f8c476') }
  ],
  crm: [
    { key: 'webhook_url', label: 'Webhook URL', placeholder: 'https://crm.example.com/api' },
    { key: 'api_key', label: 'API Key（可选）', placeholder: 'Bearer Token' }
  ]
}[source] || [])

const loadConnectors = async () => {
  try {
    const res = await http.get('/api/knowledge/connectors')
    connectors.value = res?.list || res?.data?.list || []
    for (const c of connectors.value) {
      const cfg = c.config || {}
      const form = { ...emptyCred(), enabled: c.enabled }
      for (const k of Object.keys(cfg)) {
        const v = cfg[k]
        if (typeof v === 'string' && v.startsWith('****')) {
          maskedFields.value[`${c.source}.${k}`] = true
        } else {
          form[k] = v
        }
      }
      credForm.value[c.source] = form
    }
  } catch (e) {
    console.error('load connectors failed', e)
  }
}

const maskedFields = ref({})

const isMasked = (source, key) => !!maskedFields.value[`${source}.${key}`]

const loadKBs = async () => {
  try {
    const res = await http.get('/api/knowledge-bases')
    kbList.value = res?.list || res?.data?.list || res?.data || []
  } catch (e) { console.error('load kbs failed', e) }
}

const saveCred = async (source) => {
  try {
    const payload = { ...credForm.value[source] };
    const cfg = { ...(payload.config || {}) }
    for (const [k, v] of Object.entries(cfg)) {
      if (typeof v === 'string' && v.trim() === '') delete cfg[k]
    }
    for (const f of credFields(source)) {
      if (String(payload[f.key] || '').trim() === '' && maskedFields.value[`${source}.${f.key}`]) continue
      if (f.key !== 'enabled') cfg[f.key] = payload[f.key]
    }
    await http.put(`/api/knowledge/connectors/${source}`, { enabled: payload.enabled, config: cfg })
    ElMessage.success(t('auto.Connectors.1a2c3f24'))
    await loadConnectors()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Connectors.6de920b4'))
  }
}

const testConn = async (source) => {
  try {
    const res = await http.post(`/api/knowledge/connectors/${source}/test`)
    // 拦截器已拆包：test 直接回 {ok, message, latency_ms}，按信封取 .data 只会读到 undefined
    const d = res?.data ?? res ?? {}
    d.ok ? ElMessage.success(d.message || t('auto.Connectors.b33199d3')) : ElMessage.warning(d.message || t('auto.Connectors.0745fc09'))
    await loadConnectors()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Connectors.9710d9b9'))
  }
}

const pull = async (source) => {
  const kbId = pullKb.value[source]
  if (!kbId) { ElMessage.warning(t('auto.Connectors.a56e1b8d')); return }
  pulling.value = source
  try {
    const res = await http.post(`/api/knowledge/connectors/${source}/pull`, { product_id: kbId })
    // 拦截器已拆包：pull 直接回 {source, imported, failed, skipped, details}
    const d = res?.data ?? res ?? null
    lastResult.value = d
    ElMessage.success(`拉取完成: 成功 ${d?.imported || 0} / 失败 ${d?.failed || 0}`)
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Connectors.8796b774'))
  } finally {
    pulling.value = ''
  }
}

onMounted(() => { loadConnectors(); loadKBs() })
</script>

<style scoped>
.connectors-page { padding: 20px; }
.connector-card { margin-bottom: 16px; }
.card-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.conn-name { font-weight: 600; font-size: 15px; }
.conn-desc { margin-bottom: 8px; }
.test-msg { margin-left: 6px; font-size: 12px; color: #909399; }
.form-actions { margin-top: 8px; display: flex; gap: 8px; }
.pull-box { display: flex; gap: 8px; margin-top: 8px; }
.kb-select { flex: 1; }
.result-card { margin-top: 8px; }
</style>
