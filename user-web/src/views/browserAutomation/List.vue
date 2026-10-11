<template>
  <div class="page">
    <div class="toolbar">
      <el-button type="primary" icon="Plus" @click="$router.push('/browser-automation/tasks/create')">{{ $t('auto.List.8.b255a2') }}</el-button>
      <el-select v-model="query.status" :placeholder="$t('auto.List.1.3fea7c')" clearable style="width: 140px" @change="onFilterChange">
        <el-option v-for="s in statusOptions" :key="s.value" :label="s.label" :value="s.value" />
      </el-select>
    </div>

    <el-table :data="list" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="name" :label="$t('auto.List.2.d7ec2d')" min-width="160" />
      <el-table-column prop="task_type" :label="$t('auto.List.3.226b09')" width="100" />
      <el-table-column prop="status" :label="$t('auto.List.1.3fea7c')" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="brain_mode" :label="$t('auto.List.4.f0789e')" width="90">
        <template #default="{ row }">{{ row.brain_mode ? 'Brain' : '步骤' }}</template>
      </el-table-column>
      <el-table-column prop="last_run_at" :label="$t('auto.List.5.f3d8ee')" width="170">
        <template #default="{ row }">{{ row.last_run_at ? formatTime(row.last_run_at) : '—' }}</template>
      </el-table-column>
      <el-table-column :label="$t('auto.List.6.2b6bc0')" width="330" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="success" :disabled="row.status === 'running'" @click="onRun(row)">{{ $t('auto.List.9.1a6aa2') }}</el-button>
          <el-button size="small" v-if="row.status === 'draft'" @click="onPublish(row)">{{ $t('auto.List.10.83611a') }}</el-button>
          <el-button size="small" v-if="row.status === 'running'" @click="onPause(row)">{{ $t('auto.List.11.8d63ef') }}</el-button>
          <el-button size="small" @click="$router.push(`/browser-automation/tasks/${row.id}`)">{{ $t('auto.List.12.f26225') }}</el-button>
          <el-button size="small" @click="$router.push(`/browser-automation/tasks/${row.id}/edit`)">{{ $t('auto.List.13.95b351') }}</el-button>
          <el-button size="small" type="danger" @click="onDelete(row)">{{ $t('auto.List.14.2f4aad') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="query.page"
      :page-size="query.limit"
      :total="total"
      layout="total, prev, pager, next"
      style="margin-top: 12px"
      @current-change="load"
    />

    <!-- 无 Host 引导弹窗（8001）：与详情页、Host 状态页同一份步骤 -->
    <el-dialog v-model="hostDialog.visible" :title="$t('auto.List.7.604cad')" width="520px">
      <HostInstallGuide />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  listBrowserTasks, publishBrowserTask, runBrowserTask, pauseBrowserTask, deleteBrowserTask,
} from '@/api/browserAutomation'
import { classifyRunError, HOST_OFFLINE, TASK_BUSY, RUN_ERROR_TEXT } from './hostRunError'
import HostInstallGuide from './HostInstallGuide.vue'

const router = useRouter()

const statusOptions = [
  { value: 'draft', label: t('auto.List.15.22b433') }, { value: 'ready', label: t('auto.List.16.c0d218') },
  { value: 'running', label: t('auto.List.17.46e386') }, { value: 'paused', label: t('auto.List.18.a2d930') },
  { value: 'done', label: t('auto.List.19.769d88') }, { value: 'failed', label: t('auto.List.20.acd5cb') },
  { value: 'archived', label: t('auto.List.21.c3ba16') },
]

const { t } = useI18n()

const list = ref([])
const total = ref(0)
const loading = ref(false)
const query = reactive({ status: '', page: 1, limit: 20 })
const hostDialog = reactive({ visible: false })
let pollTimer = null

const unpack = (res) => res?.data ?? res

// 换筛选条件是一次新的查询，不是同一批数据的下一页：不回到第 1 页时，
// 筛完停在自己正在看的那一页，列表空了而分页条还留着——读起来像「这个状态下没有任务」。
// 分页条自己只在总页数变小往下夹，条件变严但总页数没变的那一档它不管。
function onFilterChange() {
  query.page = 1
  load()
}

async function load() {
  loading.value = true
  try {
    const res = await listBrowserTasks({ ...query })
    const data = unpack(res)
    list.value = data?.list || (Array.isArray(data) ? data : [])
    total.value = data?.total || 0
    // 有 running 任务时轻量轮询刷新状态
    if (list.value.some((t) => t.status === 'running')) {
      if (!pollTimer) pollTimer = setInterval(load, 5000)
    } else if (pollTimer) {
      clearInterval(pollTimer); pollTimer = null
    }
  } finally {
    loading.value = false
  }
}

const statusTagType = (s) => ({
  draft: 'info', ready: '', running: 'warning', paused: 'info',
  done: 'success', failed: 'danger', archived: 'info',
}[s] || 'info')

const formatTime = (t) => new Date(t).toLocaleString('zh-CN')

function handleRunError(err) {
  const kind = classifyRunError(err)
  if (kind === HOST_OFFLINE) hostDialog.visible = true
  else if (kind === TASK_BUSY) ElMessage.warning(RUN_ERROR_TEXT(err))
  else ElMessage.error(RUN_ERROR_TEXT(err))
}

async function onRun(row) {
  try {
    const res = await runBrowserTask(row.id)
    // 点了执行就要看得见执行——直接进这条会话的监控页（放行闸门、步骤流都在那里）
    const sid = unpack(res)?.session_id
    ElMessage.success(t('auto.List.22.2bdc57'))
    load()
    if (sid) router.push(`/browser-automation/sessions/${sid}`)
  } catch (e) { handleRunError(e) }
}

async function onPublish(row) {
  try {
    await publishBrowserTask(row.id)
    ElMessage.success(t('auto.List.23.dca0c1'))
    load()
  } catch (e) { ElMessage.error(String(e?.message || e)) }
}

async function onPause(row) {
  try {
    await pauseBrowserTask(row.id)
    ElMessage.success(t('auto.List.18.a2d930'))
    load()
  } catch (e) { ElMessage.error(String(e?.message || e)) }
}

async function onDelete(row) {
  let go = true
  try {
    await ElMessageBox.confirm(`确认删除任务「${row.name}」？`, t('auto.List.24.02d981'), { type: 'warning' })
  } catch { go = false } // 用户取消不是失败：不接住它，控制台里会多一份假异常
  if (!go) return
  try {
    await deleteBrowserTask(row.id)
    ElMessage.success(t('auto.List.25.5cc232'))
  } catch (e) {
    ElMessage.error(`删除失败：${String(e?.message || e)}`)
  }
  load()
}

onMounted(load)
onBeforeUnmount(() => { if (pollTimer) clearInterval(pollTimer) })
</script>

<style scoped>
.page { padding: 16px; }
.toolbar { display: flex; gap: 12px; margin-bottom: 12px; }
</style>
