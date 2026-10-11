<template>
  <div class="automation-hub">
    
    <el-card class="mb16">
      <template #header><b>{{ $t('auto.AutomationHub.9cc0ca2c') }}</b><el-tag size="small" style="margin-left:8px">{{ oh.enabled ? '已启用' : '未启用' }}</el-tag></template>
      <el-form label-width="120px" size="small">
        <el-form-item :label="$t('auto.AutomationHub.e0d6ffb7')">
          <el-switch v-model="oh.enabled" />
        </el-form-item>
        <el-form-item :label="$t('auto.AutomationHub.9af79c85')">
          <div v-for="(r, i) in oh.daily_ranges" :key="i" style="margin-bottom:6px">
            <el-time-select v-model="r[0]" start="00:00" step="00:30" end="23:30" style="width:110px" />
            <span style="margin:0 6px">{{ $t('auto.AutomationHub.981cbe31') }}</span>
            <el-time-select v-model="r[1]" start="00:00" step="00:30" end="23:30" style="width:110px" />
            <el-button text type="danger" size="small" @click="oh.daily_ranges.splice(i, 1)">{{ $t('auto.AutomationHub.2f4aaddd') }}</el-button>
          </div>
          <el-button size="small" @click="oh.daily_ranges.push(['09:00','18:00'])">{{ $t('auto.AutomationHub.8fb3bc2c') }}</el-button>
        </el-form-item>
        <el-form-item :label="$t('auto.AutomationHub.8e8c90ed')">
          <el-input v-model="oh.away_message" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="saveOH">{{ $t('auto.AutomationHub.2d3f52b4') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    
    <el-card class="mb16">
      <template #header><b>{{ $t('auto.AutomationHub.7b0ca810') }}</b></template>
      <div style="margin-bottom:10px;display:flex;gap:8px">
        <el-input v-model="newMacro.name" :placeholder="$t('auto.AutomationHub.2508917a')" style="width:200px" size="small" />
        <el-select v-model="newMacro.type" size="small" style="width:160px" :placeholder="$t('auto.AutomationHub.c500cfab')">
          <el-option :label="$t('auto.AutomationHub.5d8b345e')" value="add_note" />
          <el-option :label="$t('auto.AutomationHub.73915d63')" value="add_tag" />
          <el-option :label="$t('auto.AutomationHub.3269f7e1')" value="set_priority" />
          <el-option :label="$t('auto.AutomationHub.94e7055b')" value="send_message" />
          <el-option :label="$t('auto.AutomationHub.f1ae05c8')" value="close" />
        </el-select>
        <el-input v-model="newMacro.value" :placeholder="$t('auto.AutomationHub.552cc64a')" size="small" style="width:260px" />
        <el-button type="primary" size="small" @click="addMacro">{{ $t('auto.AutomationHub.a02a7377') }}</el-button>
      </div>
      <el-table :data="macros" size="small">
        <el-table-column prop="name" :label="$t('auto.AutomationHub.d7ec2d3f')" min-width="140" />
        <el-table-column :label="$t('auto.AutomationHub.7f6f7581')" min-width="300">
          <template #default="{ row }">
            <el-tag v-for="(a, i) in parseActions(row.actions)" :key="i" size="small" style="margin-right:4px">{{ a.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.AutomationHub.2b6bc0f2')" width="80">
          <template #default="{ row }">
            <el-button link type="danger" size="small" @click="delMacro(row)">{{ $t('auto.AutomationHub.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <el-card class="mb16">
      <template #header><b>{{ $t('auto.AutomationHub.084e2300') }}</b><span class="hint">{{ $t('auto.AutomationHub.2c5ccbbf') }}</span></template>
      <div style="margin-bottom:10px;display:flex;gap:8px">
        <el-input v-model="newHook.url" placeholder="https://your-app.com/webhook" size="small" style="width:320px" />
        <el-input v-model="newHook.events" placeholder="all" size="small" style="width:160px" />
        <el-button type="primary" size="small" @click="addHook">{{ $t('auto.AutomationHub.40aa7e25') }}</el-button>
      </div>
      <el-table :data="hooks" size="small">
        <el-table-column prop="url" label="URL" min-width="240" show-overflow-tooltip />
        <el-table-column prop="events" :label="$t('auto.AutomationHub.10b2761d')" width="160" />
        <el-table-column :label="$t('auto.AutomationHub.3fea7ca7')" width="80">
          <template #default="{ row }"><el-tag size="small" :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '启用' : '停用' }}</el-tag></template>
        </el-table-column>
        <el-table-column :label="$t('auto.AutomationHub.2b6bc0f2')" width="80">
          <template #default="{ row }">
            <el-button link type="danger" size="small" @click="delHook(row)">{{ $t('auto.AutomationHub.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <el-card class="mb16">
      <template #header><b>{{ $t('auto.AutomationHub.3e6e56cb') }}</b><el-button size="small" style="float:right" @click="sendNow">{{ $t('auto.AutomationHub.41446fc6') }}</el-button></template>
      <div style="margin-bottom:10px;display:flex;gap:8px">
        <el-input v-model="newSub.email" :placeholder="$t('auto.AutomationHub.c59a2a14')" size="small" style="width:260px" />
        <el-select v-model="newSub.schedule" size="small" style="width:120px">
          <el-option :label="$t('auto.AutomationHub.ce221098')" value="daily" />
          <el-option :label="$t('auto.AutomationHub.09b03b6d')" value="weekly" />
        </el-select>
        <el-button type="primary" size="small" @click="addSub">{{ $t('auto.AutomationHub.02daf71f') }}</el-button>
      </div>
      <el-table :data="subs" size="small">
        <el-table-column prop="email" :label="$t('auto.AutomationHub.3bc5e602')" min-width="200" />
        <el-table-column prop="schedule" :label="$t('auto.AutomationHub.b70a5564')" width="100" />
        <el-table-column :label="$t('auto.AutomationHub.2c6aa24c')" width="170">
          <template #default="{ row }">{{ row.last_sent ? String(row.last_sent).replace('T',' ').slice(0,16) : '—' }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.AutomationHub.2b6bc0f2')" width="80">
          <template #default="{ row }">
            <el-button link type="danger" size="small" @click="delSub(row)">{{ $t('auto.AutomationHub.b30d52d6') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <el-card>
      <template #header><b>{{ $t('auto.AutomationHub.135ea783') }}</b>
        <el-radio-group v-model="perfDays" size="small" style="float:right" @change="loadPerf">
          <el-radio-button :value="7">{{ $t('auto.AutomationHub.d716067e') }}</el-radio-button>
          <el-radio-button :value="30">{{ $t('auto.AutomationHub.853f8cf7') }}</el-radio-button>
        </el-radio-group>
      </template>
      <el-row :gutter="12" v-if="perf">
        <el-col :span="5"><div class="stat"><div class="v">{{ perf.total_sessions }}</div><div class="l">{{ $t('auto.AutomationHub.f3db9e96') }}</div></div></el-col>
        <el-col :span="5"><div class="stat"><div class="v" style="color:#10b981">{{ perf.auto_rate?.toFixed(1) }}%</div><div class="l">{{ $t('auto.AutomationHub.c3f0a12e') }}</div></div></el-col>
        <el-col :span="5"><div class="stat"><div class="v">{{ perf.llm_calls }}</div><div class="l">{{ $t('auto.AutomationHub.7e1189e7') }}</div></div></el-col>
        <el-col :span="5"><div class="stat"><div class="v">${{ perf.llm_cost?.toFixed(4) }}</div><div class="l">{{ $t('auto.AutomationHub.91442f6d') }}</div></div></el-col>
        <el-col :span="4"><div class="stat"><div class="v">{{ perf.human_handled }}</div><div class="l">{{ $t('auto.AutomationHub.2bbab23b') }}</div></div></el-col>
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
    ElMessage.success(t('auto.AutomationHub.79d36f86'))
  } catch (e) { ElMessage.error(e?.message || t('auto.AutomationHub.6de920b4')) }
}

const addMacro = async () => {
  try {
    await http.post('/api/macros', { name: newMacro.value.name, actions: [{ type: newMacro.value.type, value: newMacro.value.value }] })
    ElMessage.success(t('auto.AutomationHub.db6c0970'))
    newMacro.value = { name: '', type: 'add_note', value: '' }
    loadAll()
  } catch (e) { ElMessage.error(e?.message || t('auto.AutomationHub.a889286a')) }
}

const delMacro = async (row) => {
  try {
    await ElMessageBox.confirm(`删除宏「${row.name}」？`, t('auto.AutomationHub.e83a256e'), { type: 'warning' })
    await http.delete(`/api/macros/${row.id}`)
    loadAll()
  } catch (e) { if (e !== 'cancel' && e !== 'close') ElMessage.error(t('auto.AutomationHub.acf0664a')) }
}

const addHook = async () => {
  try {
    const res = await http.post('/api/webhook-subscriptions', newHook.value)
    ElMessageBox.alert(`订阅创建成功。签名密钥（仅显示一次，用于校验推送）：\n${res?.data?.secret || res?.secret || ''}`, t('auto.AutomationHub.521392f3'))
    loadAll()
  } catch (e) { ElMessage.error(e?.message || t('auto.AutomationHub.a889286a')) }
}

const delHook = async (row) => {
  try {
    await ElMessageBox.confirm(`删除订阅 ${row.url}？`, t('auto.AutomationHub.e83a256e'), { type: 'warning' })
    await http.delete(`/api/webhook-subscriptions/${row.id}`)
    loadAll()
  } catch (e) { if (e !== 'cancel' && e !== 'close') ElMessage.error(t('auto.AutomationHub.acf0664a')) }
}

const addSub = async () => {
  try {
    await http.post('/api/report-subscriptions', newSub.value)
    ElMessage.success(t('auto.AutomationHub.e420f494'))
    newSub.value = { email: '', schedule: 'daily' }
    loadAll()
  } catch (e) { ElMessage.error(e?.message || t('auto.AutomationHub.3a89ae1c')) }
}

const delSub = async (row) => {
  try {
    await http.delete(`/api/report-subscriptions/${row.id}`)
    loadAll()
  } catch (e) { ElMessage.error(t('auto.AutomationHub.bb746798')) }
}

const sendNow = async () => {
  try {
    const res = await http.post('/api/report-subscriptions/send-now')
    ElMessage.success(`已发送 ${res?.data?.sent ?? res?.sent ?? 0} 份（未配置 SMTP 时会失败并记录日志）`)
  } catch (e) { ElMessage.error(e?.message || t('auto.AutomationHub.9ca6a344')) }
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
