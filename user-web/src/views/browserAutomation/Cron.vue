<template>
  <div class="page">
    <div class="toolbar">
      <h2 style="margin: 0">{{ $t('auto.Cron.00219842') }}</h2>
      <el-space>
        <el-button type="primary" icon="Plus" @click="openCreate">{{ $t('auto.Cron.1aeb9ead') }}</el-button>
        <el-button @click="$router.push('/browser-automation/tasks/create')">{{ $t('auto.Cron.d66eefc7') }}</el-button>
      </el-space>
    </div>

    <el-table :data="list" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="task_id" :label="$t('auto.Cron.0e46d8d6')" min-width="230">
        <template #default="{ row }">
          <el-link type="primary" @click="$router.push(`/browser-automation/tasks/${row.task_id}`)">#{{ row.task_id }} {{ taskName(row.task_id) }}</el-link>
          <template v-if="carrierKnown(row.task_id)">
            <el-tag size="small" :type="statusTagType(taskStatus(row.task_id))" class="task-tag">
              {{ taskStatus(row.task_id) }}
            </el-tag>
            <div v-if="row.enabled && !taskRunnable(row.task_id)" class="task-warn">
              任务状态 {{ taskStatus(row.task_id) }} 不可执行，到点只会空跑（不产生会话）——需先发布
            </div>
          </template>
          <template v-else>
            <el-tag type="danger" size="small" class="task-tag">{{ $t('auto.Cron.b563b577') }}</el-tag>
            <div class="task-warn">
              这条触发器的任务已不在任务列表里（可能已被删除）。它到点执行不了，请删除本条
            </div>
          </template>
        </template>
      </el-table-column>
      <el-table-column prop="cron_expr" :label="$t('auto.Cron.a9400c40')" width="160" />
      <el-table-column prop="time_zone" :label="$t('auto.Cron.1220b05a')" width="140">
        <template #default="{ row }">{{ row.time_zone || '服务器本地' }}</template>
      </el-table-column>
      <el-table-column prop="enabled" :label="$t('auto.Cron.3fea7ca7')" width="100">
        <template #default="{ row }">
          <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '已启用' : '已停用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="next_run_at" :label="$t('auto.Cron.0383450d')" width="220">
        <template #default="{ row }">
          <template v-if="!row.enabled">{{ $t('auto.Cron.86e42213') }}</template>
          <template v-else-if="!row.next_run_at">
            <el-tag type="danger" size="small">{{ $t('auto.Cron.80463788') }}</el-tag>
            <span class="next-hint">{{ $t('auto.Cron.5783937a') }}</span>
          </template>
          <template v-else>
            {{ formatTime(row.next_run_at) }}
            <span class="next-hint">{{ relativeToNow(row.next_run_at) }}</span>
          </template>
        </template>
      </el-table-column>
      <el-table-column prop="last_run_at" :label="$t('auto.Cron.a1b7bec0')" width="170">
        <template #default="{ row }">{{ row.last_run_at ? formatTime(row.last_run_at) : '—' }}</template>
      </el-table-column>
      <el-table-column :label="$t('auto.Cron.2b6bc0f2')" width="280" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="toggle(row)">{{ row.enabled ? '停用' : '启用' }}</el-button>
          <el-button size="small" @click="openEdit(row)">{{ $t('auto.Cron.95b351c8') }}</el-button>
          <el-button size="small" type="danger" @click="remove(row)">{{ $t('auto.Cron.2f4aaddd') }}</el-button>
        </template>
      </el-table-column>
      <template #empty>
        <div class="empty-hint">
          还没有触发器。触发器只能挂在 <b>{{ $t('auto.Cron.9bb577cd') }}</b> 的任务上，且一条任务最多一个。
        </div>
      </template>
    </el-table>

    <el-dialog v-model="dialog.visible" ::title="$t('auto.Cron.0213ddf5')" width="520px">
      <el-form :model="dialog" label-width="100px">
        <el-form-item :label="$t('auto.Cron.0e46d8d6')" required>
          <el-select v-model="dialog.task_id" :disabled="!!dialog.id" :placeholder="$t('auto.Cron.b3d2d54f')" style="width: 100%">
            <el-option v-for="t in cronTasks" :key="t.id" :label="`#${t.id} ${t.name}`" :value="t.id" />
          </el-select>
          <div v-if="!dialog.id && !cronTasks.length" class="form-hint">
            没有可选任务：先把某条任务的类型改成「定时 (cron)」并发布
          </div>
          <div v-else-if="!dialog.id" class="form-hint">{{ $t('auto.Cron.b866f3ec') }}</div>
          <div v-if="dialog.task_id && !taskRunnable(dialog.task_id)" class="form-warn">
            已选任务状态为 {{ taskStatus(dialog.task_id) }}：触发器会照常排程，但到点会被执行入口拒掉、不产生会话——请先发布该任务
          </div>
        </el-form-item>
        <el-form-item :label="$t('auto.Cron.d8178700')" required>
          <el-input v-model="dialog.cron_expr" :placeholder="$t('auto.Cron.12182670')" />
          <div class="form-hint">{{ $t('auto.Cron.a5923e90') }}</div>
        </el-form-item>
        <el-form-item :label="$t('auto.Cron.1220b05a')">
          <el-select v-model="dialog.time_zone" clearable :placeholder="$t('auto.Cron.1abc489f')" style="width: 100%">
            <el-option v-for="tz in TIMEZONES" :key="tz" :label="tz" :value="tz" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="!dialog.id" :label="$t('auto.Cron.1388f233')">
          <el-switch v-model="dialog.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialog.visible = false">{{ $t('auto.Cron.625fb26b') }}</el-button>
        <el-button type="primary" :loading="dialog.saving" @click="submit">{{ $t('auto.Cron.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  listBrowserCron, createBrowserCron, updateBrowserCron,
  enableBrowserCron, disableBrowserCron, deleteBrowserCron,
  listBrowserTasks,
} from '@/api/browserAutomation'

const TIMEZONES = ['Asia/Shanghai', 'Asia/Tokyo', 'Europe/London', 'America/New_York', 'UTC']

const { t } = useI18n()

const list = ref([])
const tasks = ref([])
const loading = ref(false)
const dialog = ref({ visible: false, saving: false, id: null, task_id: null, cron_expr: '', time_zone: '', enabled: true })

const unpack = (res) => res?.data ?? res
const asList = (data) => (Array.isArray(data) ? data : data?.list || [])

const taskName = (id) => tasks.value.find((t) => t.id === id)?.name || ''
const taskStatus = (id) => tasks.value.find((t) => t.id === id)?.status || ''
const carrierKnown = (id) => tasks.value.some((t) => t.id === id)

// 服务端只对 task_type=cron 校验，不校验状态：排程照挂、next_run_at 照推进，
// 但真正执行时会被「任务状态 X 不可执行（需先 publish）」拒掉——于是触发器看起来在跑，
// 库里却一条会话都不产生。可执行状态口径与服务端执行入口保持一致。
const RUNNABLE_STATUS = ['ready', 'paused', 'done', 'failed']
const taskRunnable = (id) => RUNNABLE_STATUS.includes(taskStatus(id))
const statusTagType = (s) => ({
  draft: 'info', ready: '', running: 'warning', paused: 'info',
  done: 'success', failed: 'danger', archived: 'info',
}[s] ?? 'info')

// 触发器只认 cron 任务；已占用那条任务的选择项要摘掉——服务端对 task_id 建了唯一索引，
// 留着选项等于每次都要撞一次 409 才知道不行。
const cronTasks = computed(() => {
  const taken = new Set(list.value.map((c) => c.task_id))
  return tasks.value.filter((t) => t.task_type === 'cron' && (!taken.has(t.id) || t.id === dialog.value.task_id))
})

const formatTime = (v) => {
  if (!v) return '—'
  const d = new Date(v)
  return Number.isNaN(d.getTime()) ? String(v) : d.toLocaleString('zh-CN')
}

const relativeToNow = (v) => {
  const ms = new Date(v).getTime() - Date.now()
  if (Number.isNaN(ms)) return ''
  if (ms <= 0) return t('auto.Cron.ec01ae45')
  const min = Math.round(ms / 60000)
  if (min < 60) return `约 ${min} 分钟后`
  const hour = Math.round(min / 60)
  if (hour < 24) return `约 ${hour} 小时后`
  return `约 ${Math.round(hour / 24)} 天后`
}

async function load() {
  loading.value = true
  try {
    const [cRes, tRes] = await Promise.all([
      listBrowserCron(),
      listBrowserTasks({ limit: 200 }),
    ])
    list.value = asList(unpack(cRes))
    tasks.value = asList(unpack(tRes))
  } finally {
    loading.value = false
  }
}

function openCreate() {
  dialog.value = { visible: true, saving: false, id: null, task_id: null, cron_expr: '*/30 * * * *', time_zone: 'Asia/Shanghai', enabled: true }
}

function openEdit(row) {
  dialog.value = { visible: true, saving: false, id: row.id, task_id: row.task_id, cron_expr: row.cron_expr, time_zone: row.time_zone || '', enabled: row.enabled }
}

async function submit() {
  const d = dialog.value
  if (!d.cron_expr.trim() || (!d.id && !d.task_id)) {
    ElMessage.warning(t('auto.Cron.3318d683'))
    return
  }
  d.saving = true
  try {
    if (d.id) await updateBrowserCron(d.id, { cron_expr: d.cron_expr.trim(), time_zone: d.time_zone || '' })
    else await createBrowserCron({ task_id: d.task_id, cron_expr: d.cron_expr.trim(), time_zone: d.time_zone || '', enabled: d.enabled })
    ElMessage.success(d.id ? t('auto.Cron.38efe72e') : t('auto.Cron.698ad53c'))
    d.visible = false
  } catch {
    // 拒因（表达式第几段不对、该任务已有触发器）由服务端文案给出，拦截器已经弹过一次；
    // 这里不重复弹，只保证列表回到与后端一致的状态。
  } finally {
    d.saving = false
    load()
  }
}

async function toggle(row) {
  try {
    if (row.enabled) await disableBrowserCron(row.id)
    else await enableBrowserCron(row.id)
    ElMessage.success(t('auto.Cron.1b92092d'))
  } catch {
    // 启停失败（例如表达式这次算不出排程）必须让界面回到后端真相，
    // 不能把「点了没生效」的旧状态留着让用户再点一次。
  } finally {
    load()
  }
}

async function remove(row) {
  try {
    await ElMessageBox.confirm(`确认删除任务「${taskName(row.task_id) || '#' + row.task_id}」的触发器？删除后该任务不会再被定时唤起。`, t('auto.Cron.02d9819d'), { type: 'warning' })
  } catch {
    return // 用户取消
  }
  try {
    await deleteBrowserCron(row.id)
    ElMessage.success(t('auto.Cron.5cc23262'))
  } catch {
    // 同上：失败原因已由拦截器提示，此处只负责刷新回真相
  } finally {
    load()
  }
}

onMounted(load)
</script>

<style scoped>
.page { padding: 16px; }
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.next-hint { color: #909399; font-size: 12px; margin-left: 6px; }
.form-hint { color: #909399; font-size: 12px; line-height: 1.6; }
.form-warn { color: #e6a23c; font-size: 12px; line-height: 1.6; }
.task-tag { margin-left: 6px; }
.task-warn { color: #e6a23c; font-size: 12px; line-height: 1.5; margin-top: 2px; }
.empty-hint { color: #909399; font-size: 13px; padding: 12px; }
</style>
