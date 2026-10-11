<template>
  <div class="geo-page">
    <div class="p-4">
      
      <el-row :gutter="16" class="mb-4">
        <el-col :span="8">
          <el-card>
            <el-statistic :title="$t('auto.AlertCenter.1ba23536')" :value="unreadCount">
              <template #suffix>
                <el-tag v-if="unreadCount > 0" type="danger" size="small" effect="dark" class="ml-2">{{ $t('auto.AlertCenter.2d25e6f9') }}</el-tag>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card><el-statistic :title="$t('auto.AlertCenter.1b290dde')" :value="total" /></el-card>
        </el-col>
        <el-col :span="8">
          <el-card>
            <div class="flex flex-col gap-2">
              <el-button type="primary" :loading="loading" @click="load">{{ $t('auto.AlertCenter.694fc5ef') }}</el-button>
              <el-button @click="ackAllVisible" :disabled="!hasUnackInPage">{{ $t('auto.AlertCenter.6f172f6a') }}</el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>

      
      <el-card class="mb-4">
        <div class="flex items-center gap-3 flex-wrap">
          <el-select v-model="filterType" :placeholder="$t('auto.AlertCenter.c62e34c5')" clearable style="width:180px" @change="load">
            <el-option :label="$t('auto.AlertCenter.53119fb3')" value="negative_monitor" />
            <el-option :label="$t('auto.AlertCenter.0ae1dad6')" value="sov_drop" />
            <el-option :label="$t('auto.AlertCenter.53dadcf1')" value="entity_anomaly" />
          </el-select>
          <el-select v-model="filterLevel" :placeholder="$t('auto.AlertCenter.e78e4b2d')" clearable style="width:140px" @change="load">
            <el-option label="Info" value="info" />
            <el-option label="Warning" value="warning" />
            <el-option label="Critical" value="critical" />
          </el-select>
        </div>
      </el-card>

      
      <el-card>
        <template #header><span class="font-bold">{{ $t('auto.AlertCenter.f0a13c25') }}</span></template>
        <el-table :data="alerts" v-loading="loading" size="small">
          <el-table-column :label="$t('auto.AlertCenter.3fea7ca7')" width="90">
            <template #default="{ row }">
              <el-tag v-if="!row.notified" type="danger" size="small">{{ $t('auto.AlertCenter.fc8ee819') }}</el-tag>
              <el-tag v-else type="info" size="small">{{ $t('auto.AlertCenter.4113e726') }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.AlertCenter.e78e4b2d')" width="100">
            <template #default="{ row }">
              <el-tag :type="levelTagType(row.level)" size="small">{{ row.level }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.AlertCenter.226b0912')" width="130">
            <template #default="{ row }">
              <el-tag size="small" effect="plain">{{ typeLabel(row.type) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="engine" :label="$t('auto.AlertCenter.6c1b19dd')" width="110" />
          <el-table-column prop="query" :label="$t('auto.AlertCenter.97d095e0')" min-width="180" show-overflow-tooltip />
          <el-table-column prop="snippet" :label="$t('auto.AlertCenter.f8605e6c')" min-width="240" show-overflow-tooltip />
          <el-table-column :label="$t('auto.AlertCenter.19fcb9eb')" width="170">
            <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
          </el-table-column>
          <el-table-column :label="$t('auto.AlertCenter.2b6bc0f2')" width="150" fixed="right">
            <template #default="{ row }">
              <el-button v-if="!row.notified" size="small" type="primary" link @click="ack(row)">{{ $t('auto.AlertCenter.e83a256e') }}</el-button>
              <el-button size="small" type="danger" link @click="remove(row)">{{ $t('auto.AlertCenter.2f4aaddd') }}</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div class="flex justify-end mt-3">
          <el-pagination
            v-model:current-page="page"
            :page-size="limit"
            :total="total"
            layout="prev, pager, next, total"
            @current-change="load"
          />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listGeoAlerts, getGeoAlertsUnreadCount, ackGeoAlert, deleteGeoAlert } from '@/api/geoAlert'

const { t } = useI18n()

const alerts = ref([])
const total = ref(0)
const unreadCount = ref(0)
const loading = ref(false)
const page = ref(1)
const limit = ref(20)
const filterType = ref('')
const filterLevel = ref('')
let pollTimer = null

const hasUnackInPage = computed(() => alerts.value.some((a) => !a.notified))

const levelTagType = (level) => ({ critical: 'danger', warning: 'warning', info: 'info' }[level] || 'info')
const typeLabel = (type) => ({ negative_monitor: t('auto.AlertCenter.53119fb3'), sov_drop: 'SOV下降', entity_anomaly: t('auto.AlertCenter.53dadcf1') }[type] || type)

const formatTime = (t) => {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN', { hour12: false })
}

const load = async () => {
  loading.value = true
  try {
    const res = await listGeoAlerts({
      page: page.value,
      limit: limit.value,
      type: filterType.value || undefined,
      level: filterLevel.value || undefined
    })
    const data = res?.data || res
    alerts.value = data?.list || []
    total.value = data?.total || 0
    const unreadRes = await getGeoAlertsUnreadCount()
    const unreadData = unreadRes?.data || unreadRes
    unreadCount.value = unreadData?.count || 0
  } catch (e) {
    console.error(t('auto.AlertCenter.9f3d4139'), e)
  } finally {
    loading.value = false
  }
}

const ack = async (row) => {
  await ackGeoAlert(row.id)
  ElMessage.success(t('auto.AlertCenter.4113e726'))
  await load()
}

const ackAllVisible = async () => {
  const targets = alerts.value.filter((a) => !a.notified)
  if (!targets.length) return
  await ElMessageBox.confirm(`确认本页 ${targets.length} 条告警？`, t('auto.AlertCenter.3d3b8f94'))
  await Promise.all(targets.map((a) => ackGeoAlert(a.id)))
  ElMessage.success(`已确认 ${targets.length} 条`)
  await load()
}

const remove = async (row) => {
  await ElMessageBox.confirm(t('auto.AlertCenter.c64c7c1a'), t('auto.AlertCenter.2f4aaddd'))
  await deleteGeoAlert(row.id)
  ElMessage.success(t('auto.AlertCenter.5cc23262'))
  await load()
}

onMounted(() => {
  load()
  pollTimer = setInterval(() => {
    if (document.visibilityState === 'visible')
      load();
  }, 60000)
})
onBeforeUnmount(() => {
  if (pollTimer) clearInterval(pollTimer)
})
</script>
