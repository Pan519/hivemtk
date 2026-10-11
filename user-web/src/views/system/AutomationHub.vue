<template>
  <div class="automation-hub">
    
    <el-card class="mb16">
      <template #header><b>{{ $t('auto.AutomationHub.23.9cc0ca') }}</b><el-tag size="small" style="margin-left:8px">{{ oh.enabled ? '已启用' : '未启用' }}</el-tag></template>
      <el-form label-width="120px" size="small">
        <el-form-item :label="$t('auto.AutomationHub.1.e0d6ff')">
          <el-switch v-model="oh.enabled" />
        </el-form-item>
        <el-form-item :label="$t('auto.AutomationHub.2.9af79c')">
          <div v-for="(r, i) in oh.daily_ranges" :key="i" style="margin-bottom:6px">
            <el-time-select v-model="r[0]" start="00:00" step="00:30" end="23:30" style="width:110px" />
            <span style="margin:0 6px">{{ $t('auto.AutomationHub.24.981cbe') }}</span>
            <el-time-select v-model="r[1]" start="00:00" step="00:30" end="23:30" style="width:110px" />
            <el-button text type="danger" size="small" @click="oh.daily_ranges.splice(i, 1)">{{ $t('auto.AutomationHub.25.2f4aad') }}</el-button>
          </div>
          <el-button size="small" @click="oh.daily_ranges.push(['09:00','18:00'])">{{ $t('auto.AutomationHub.26.8fb3bc') }}</el-button>
        </el-form-item>
        <el-form-item :label="$t('auto.AutomationHub.3.8e8c90')">
          <el-input v-model="oh.away_message" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="saveOH">{{ $t('auto.AutomationHub.27.2d3f52') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    
    <el-card class="mb16">
      <template #header><b>{{ $t('auto.AutomationHub.28.7b0ca8') }}</b></template>
      <div style="margin-bottom:10px;display:flex;gap:8px">
        <el-input v-model="newMacro.name" :placeholder="$t('auto.AutomationHub.4.250891')" style="width:200px" size="small" />
        <el-select v-model="newMacro.type" size="small" style="width:160px" :placeholder="$t('auto.AutomationHub.5.c500cf')">
          <el-option :label="$t('auto.AutomationHub.6.5d8b34')" value="add_note" />
          <el-option :label="$t('auto.AutomationHub.7.73915d')" value="add_tag" />
          <el-option :label="$t('auto.AutomationHub.8.3269f7')" value="set_priority" />
          <el-option :label="$t('auto.AutomationHub.9.94e705')" value="send_message" />
          <el-option :label="$t('auto.AutomationHub.10.f1ae05')" value="close" />
        </el-select>
        <el-input v-model="newMacro.value" :placeholder="$t('auto.AutomationHub.11.552cc6')" size="small" style="width:260px" />
        <el-button type="primary" size="small" @click="addMacro">{{ $t('auto.AutomationHub.29.a02a73') }}</el-button>
      </div>
      <el-table :data="macros" size="small">
        <el-table-column prop="name" :label="$t('auto.AutomationHub.12.d7ec2d')" min-width="140" />
        <el-table-column :label="$t('auto.AutomationHub.13.7f6f75')" min-width="300">
          <template #default="{ row }">
            <el-tag v-for="(a, i) in parseActions(row.actions)" :key="i" size="small" style="margin-right:4px">{{ a.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.AutomationHub.14.2b6bc0')" width="80">
          <template #default="{ row }">
            <el-button link type="danger" size="small" @click="delMacro(row)">{{ $t('auto.AutomationHub.25.2f4aad') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <el-card class="mb16">
      <template #header><b>{{ $t('auto.AutomationHub.30.084e23') }}</b><span class="hint">{{ $t('auto.AutomationHub.31.2c5ccb') }}</span></template>
      <div style="margin-bottom:10px;display:flex;gap:8px">
        <el-input v-model="newHook.url" placeholder="https://your-app.com/webhook" size="small" style="width:320px" />
        <el-input v-model="newHook.events" placeholder="all" size="small" style="width:160px" />
        <el-button type="primary" size="small" @click="addHook">{{ $t('auto.AutomationHub.32.40aa7e') }}</el-button>
      </div>
      <el-table :data="hooks" size="small">
        <el-table-column prop="url" label="URL" min-width="240" show-overflow-tooltip />
        <el-table-column prop="events" :label="$t('auto.AutomationHub.15.10b276')" width="160" />
        <el-table-column :label="$t('auto.AutomationHub.16.3fea7c')" width="80">
          <template #default="{ row }"><el-tag size="small" :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '启用' : '停用' }}</el-tag></template>
        </el-table-column>
        <el-table-column :label="$t('auto.AutomationHub.14.2b6bc0')" width="80">
          <template #default="{ row }">
            <el-button link type="danger" size="small" @click="delHook(row)">{{ $t('auto.AutomationHub.25.2f4aad') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <el-card class="mb16">
      <template #header><b>{{ $t('auto.AutomationHub.33.3e6e56') }}</b><el-button size="small" style="float:right" @click="sendNow">{{ $t('auto.AutomationHub.34.41446f') }}</el-button></template>
      <div style="margin-bottom:10px;display:flex;gap:8px">
        <el-input v-model="newSub.email" :placeholder="$t('auto.AutomationHub.17.c59a2a')" size="small" style="width:260px" />
        <el-select v-model="newSub.schedule" size="small" style="width:120px">
          <el-option :label="$t('auto.AutomationHub.18.ce2210')" value="daily" />
          <el-option :label="$t('auto.AutomationHub.19.09b03b')" value="weekly" />
        </el-select>
        <el-button type="primary" size="small" @click="addSub">{{ $t('auto.AutomationHub.35.02daf7') }}</el-button>
      </div>
      <el-table :data="subs" size="small">
        <el-table-column prop="email" :label="$t('auto.AutomationHub.20.3bc5e6')" min-width="200" />
        <el-table-column prop="schedule" :label="$t('auto.AutomationHub.21.b70a55')" width="100" />
        <el-table-column :label="$t('auto.AutomationHub.22.2c6aa2')" width="170">
          <template #default="{ row }">{{ row.last_sent ? String(row.last_sent).replace('T',' ').slice(0,16) : '—' }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.AutomationHub.14.2b6bc0')" width="80">
          <template #default="{ row }">
            <el-button link type="danger" size="small" @click="delSub(row)">{{ $t('auto.AutomationHub.36.b30d52') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <el-card>
      <template #header><b>{{ $t('auto.AutomationHub.37.135ea7') }}</b>
        <el-radio-group v-model="perfDays" size="small" style="float:right" @change="loadPerf">
          <el-radio-button :value="7">{{ $t('auto.AutomationHub.38.d71606') }}</el-radio-button>
          <el-radio-button :value="30">{{ $t('auto.AutomationHub.39.853f8c') }}</el-radio-button>
        </el-radio-group>
      </template>
      <el-row :gutter="12" v-if="perf">
        <el-col :span="5"><div class="stat"><div class="v">{{ perf.total_sessions }}</div><div class="l">{{ $t('auto.AutomationHub.40.f3db9e') }}</div></div></el-col>
        <el-col :span="5"><div class="stat"><div class="v" style="color:#10b981">{{ perf.auto_rate?.toFixed(1) }}%</div><div class="l">{{ $t('auto.AutomationHub.41.c3f0a1') }}</div></div></el-col>
        <el-col :span="5"><div class="stat"><div class="v">{{ perf.llm_calls }}</div><div class="l">{{ $t('auto.AutomationHub.42.7e1189') }}</div></div></el-col>
        <el-col :span="5"><div class="stat"><div class="v">${{ perf.llm_cost?.toFixed(4) }}</div><div class="l">{{ $t('auto.AutomationHub.43.91442f') }}</div></div></el-col>
        <el-col :span="4"><div class="stat"><div class="v">{{ perf.human_handled }}</div><div class="l">{{ $t('auto.AutomationHub.44.2bbab2') }}</div></div></el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n';
import { ElMessage, ElMessageBox } from 'element-plus'
import { http } from '@/utils/request'

const { t } = useI18n()

const oh = ref({ enabled: false, daily_ranges: [['09:00', '18:00']], away_message: '' })
const macros = ref([])
const hooks = ref([])
const subs = ref([])
const perf = ref(null)
const perfDays = ref(7)
const newMacro = ref({ name: '', type: 'add_note', value: '' })
const newHook = ref({ url: '', events: 'all' })
const newSub = ref({ email: '', schedule: 'daily' })

const parseActions = (raw) => { try { return JSON.parse(raw || '[]') } catch { return [] } }

const loadAll = async () => {
  const [ohRes, macroRes, hookRes, subRes] = await Promise.all([
    http.get('/api/office-hours').catch(() => null),
    http.get('/api/macros').catch(() => null),
    http.get('/api/webhook-subscriptions').catch(() => null),
    http.get('/api/report-subscriptions').catch(() => null)
  ])
  // 响应拦截器已把 {code,data} 拆包：/api/office-hours 直接回配置对象，
  // 列表端点直接回 {list,total}。按信封再取 .data 只会读到 undefined。
  const ohData = ohRes?.data ?? ohRes
  if (ohData) oh.value = { ...oh.value, ...ohData }
  macros.value = macroRes?.list ?? []
  hooks.value = hookRes?.list ?? []
  subs.value = subRes?.list ?? []
  loadPerf()
}

const loadPerf = async () => {
  const res = await http.get(`/api/analytics/ai-performance?days=${perfDays.value}`).catch(() => null)
  perf.value = res?.data ?? res ?? null
}

const saveOH = async () => {
  try {
    await http.put('/api/office-hours', oh.value)
    ElMessage.success(t('auto.AutomationHub.45.79d36f'))
  } catch (e) { ElMessage.error(e?.message || t('auto.AutomationHub.46.6de920')) }
}

const addMacro = async () => {
  try {
    await http.post('/api/macros', { name: newMacro.value.name, actions: [{ type: newMacro.value.type, value: newMacro.value.value }] })
    ElMessage.success(t('auto.AutomationHub.47.db6c09'))
    newMacro.value = { name: '', type: 'add_note', value: '' }
    loadAll()
  } catch (e) { ElMessage.error(e?.message || t('auto.AutomationHub.48.a88928')) }
}

const delMacro = async (row) => {
  try {
    await ElMessageBox.confirm(`删除宏「${row.name}」？`, t('auto.AutomationHub.49.e83a25'), { type: 'warning' })
    await http.delete(`/api/macros/${row.id}`)
    loadAll()
  } catch (e) { if (e !== 'cancel' && e !== 'close') ElMessage.error(t('auto.AutomationHub.50.acf066')) }
}

const addHook = async () => {
  try {
    const res = await http.post('/api/webhook-subscriptions', newHook.value)
    ElMessageBox.alert(`订阅创建成功。签名密钥（仅显示一次，用于校验推送）：\n${res?.data?.secret || res?.secret || ''}`, t('auto.AutomationHub.51.521392'))
    loadAll()
  } catch (e) { ElMessage.error(e?.message || t('auto.AutomationHub.48.a88928')) }
}

const delHook = async (row) => {
  try {
    await ElMessageBox.confirm(`删除订阅 ${row.url}？`, t('auto.AutomationHub.49.e83a25'), { type: 'warning' })
    await http.delete(`/api/webhook-subscriptions/${row.id}`)
    loadAll()
  } catch (e) { if (e !== 'cancel' && e !== 'close') ElMessage.error(t('auto.AutomationHub.50.acf066')) }
}

const addSub = async () => {
  try {
    await http.post('/api/report-subscriptions', newSub.value)
    ElMessage.success(t('auto.AutomationHub.52.e420f4'))
    newSub.value = { email: '', schedule: 'daily' }
    loadAll()
  } catch (e) { ElMessage.error(e?.message || t('auto.AutomationHub.53.3a89ae')) }
}

const delSub = async (row) => {
  try {
    await http.delete(`/api/report-subscriptions/${row.id}`)
    loadAll()
  } catch (e) { ElMessage.error(t('auto.AutomationHub.54.bb7467')) }
}

const sendNow = async () => {
  try {
    const res = await http.post('/api/report-subscriptions/send-now')
    ElMessage.success(`已发送 ${res?.data?.sent ?? res?.sent ?? 0} 份（未配置 SMTP 时会失败并记录日志）`)
  } catch (e) { ElMessage.error(e?.message || t('auto.AutomationHub.55.9ca6a3')) }
}

onMounted(loadAll)
</script>

<style scoped>
.automation-hub { padding: 16px; }
.mb16 { margin-bottom: 16px; }
.hint { font-size: 12px; color: #94a3b8; margin-left: 10px; }
.stat { text-align: center; padding: 10px 0; background: #f8fafc; border-radius: 8px; }
.stat .v { font-size: 24px; font-weight: 700; }
.stat .l { font-size: 12px; color: #64748b; margin-top: 4px; }
</style>
