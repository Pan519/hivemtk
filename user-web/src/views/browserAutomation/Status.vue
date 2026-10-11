<template>
  <div class="page">
    <h2>{{ $t('auto.Status.8.18da72') }}</h2>

    <el-card style="margin-top: 12px">
      <el-space direction="vertical" alignment="flex-start" :size="12" style="width: 100%">
        <el-space>
          <el-button :loading="loading" @click="load">{{ $t('auto.Status.9.694fc5') }}</el-button>
          <el-button type="warning" @click="onResetToken">{{ $t('auto.Status.10.d206bc') }}</el-button>
        </el-space>

        <el-alert v-if="!loading && hosts.length === 0" type="warning" :closable="false" show-icon
          :title="$t('auto.Status.1.b2cc27')">
          <p>{{ $t('auto.Status.11.2205da') }}</p>
          <HostInstallGuide />
        </el-alert>

        <div v-if="hosts.length" class="summary">
          连接 {{ summary.count }} 台 · 可服务 {{ summary.servableCount }} 台 · 最近成功回包 {{ summary.lastOkText }}
        </div>

        <el-table v-if="hosts.length" :data="hosts" border>
          <el-table-column prop="user_id" :label="$t('auto.Status.2.e53376')" width="120" />
          <el-table-column prop="version" :label="$t('auto.Status.3.fc2aac')" width="120" />
          <el-table-column prop="pid" :label="$t('auto.Status.4.841d1a')" width="120" />
          <el-table-column :label="$t('auto.Status.5.3fea7c')" width="130">
            <template #default="{ row }">
              <el-tag v-if="row.online && row.servable" type="success">{{ $t('auto.Status.12.d2d888') }}</el-tag>
              <el-tag v-else-if="row.online" type="danger">{{ $t('auto.Status.13.6544ae') }}</el-tag>
              <el-tag v-else type="info">{{ $t('auto.Status.14.3842ba') }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.Status.6.b2f3fe')" min-width="180">
            <template #default="{ row }">
              <span v-if="row.last_cmd_ok_at">{{ formatTime(row.last_cmd_ok_at) }}</span>
              <el-text v-else type="warning">{{ $t('auto.Status.15.6b5c63') }}</el-text>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.Status.7.f1e31a')" min-width="220">
            <template #default="{ row }">
              <span v-if="!row.online">{{ $t('auto.Status.16.055404') }}</span>
              <span v-else-if="!row.servable">{{ $t('auto.Status.17.8f71e7') }}</span>
              <span v-else>{{ $t('auto.Status.18.80364e') }}</span>
            </template>
          </el-table-column>
        </el-table>
      </el-space>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getBrowserHostStatus, resetBrowserHostToken } from '@/api/browserAutomation'
import HostInstallGuide from './HostInstallGuide.vue'

const { t } = useI18n()

const hosts = ref([])
const loading = ref(false)
const summaryRaw = ref({ count: 0, servable: false, last_cmd_ok_at: null })

const unpack = (res) => res?.data ?? res
const formatTime = (t) => new Date(t).toLocaleString('zh-CN')

// 顶层 count/servable 由服务端汇总（普通用户只看自己的连接），列表按 row 出读数，
// 两者取其一即可，不一致时以行数为准——row 才是「哪一台不能干活」的答案。
const summary = computed(() => {
  const rows = hosts.value
  return {
    count: rows.length || summaryRaw.value.count || 0,
    servableCount: rows.filter((h) => h.online && h.servable).length,
    lastOkText: (() => {
      const latest = rows.map((h) => h.last_cmd_ok_at).filter(Boolean).sort().pop()
        || summaryRaw.value.last_cmd_ok_at
      return latest ? formatTime(latest) : '—'
    })(),
  }
})

async function load() {
  loading.value = true
  try {
    const res = await getBrowserHostStatus()
    const data = unpack(res)
    hosts.value = data?.hosts || []
    summaryRaw.value = { count: data?.count ?? 0, servable: data?.servable === true, last_cmd_ok_at: data?.last_cmd_ok_at ?? null }
  } catch (e) {
    hosts.value = []
    ElMessage.error(String(e?.message || e))
  } finally {
    loading.value = false
  }
}

async function onResetToken() {
  let go = true
  try {
    await ElMessageBox.confirm(
      t('auto.Status.19.6c62b6'),
      t('auto.Status.10.d206bc'), { type: 'warning' },
    )
  } catch { go = false } // 用户取消不是失败
  if (!go) return
  try {
    const res = await resetBrowserHostToken()
    const data = unpack(res)
    await ElMessageBox.alert(
      `新 token：<code style="user-select:all">${data?.token || ''}</code><br/>请将其填入本机 ~/.hivemtk/nm_host.conf 的 token= 行，然后重启 Host 进程`,
      t('auto.Status.20.0642b7'), { dangerouslyUseHTMLString: true },
    )
  } catch (e) {
    ElMessage.error(`重置失败：${String(e?.message || e)}`)
  }
}

onMounted(load)
</script>

<style scoped>
.page { padding: 16px; }
.summary { color: #606266; font-size: 13px; }
</style>
