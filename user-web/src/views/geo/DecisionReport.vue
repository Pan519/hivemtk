<template>
  <div class="geo-page">

    <div class="p-4">
    <el-row :gutter="16" class="mb-4">
      <el-col :span="8"><el-card><el-statistic :title="$t('auto.DecisionReport.6ff099c0')" :value="report.leads_captured || 0" /></el-card></el-col>
      <el-col :span="8"><el-card><el-statistic :title="$t('auto.DecisionReport.460225f6')" :value="report.tasks_pending || 0" /></el-card></el-col>
      <el-col :span="8"><el-card><el-statistic :title="$t('auto.DecisionReport.08fd5738')" :value="report.tasks_done || 0" /></el-card></el-col>
    </el-row>

    <el-card>
      <template #header><span>{{ $t('auto.DecisionReport.397bda69') }}</span></template>
      <el-table :data="tasks" v-loading="loading" size="small">
        <el-table-column prop="keyword" :label="$t('auto.DecisionReport.9699a50e')" min-width="140" />
        <el-table-column prop="intent" :label="$t('auto.DecisionReport.764aca4e')" width="90" />
        <el-table-column prop="gap_type" :label="$t('auto.DecisionReport.e755cd68')" width="150" />
        <el-table-column prop="detail" :label="$t('auto.DecisionReport.f26225bd')" min-width="260" show-overflow-tooltip />
        <el-table-column :label="$t('auto.DecisionReport.2b6bc0f2')" width="100">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="onDone(row)">{{ $t('auto.DecisionReport.769d88e4') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { getGeoDecisionReport, getGeoGapTasks, markGeoTaskDone } from '@/api/geoDecision.js'
import { toList } from '@/utils/list'

const { t } = useI18n()

const report = ref({})
const tasks = ref([])
const loading = ref(false)

const load = async () => {
  loading.value = true
  try {
    report.value = (await getGeoDecisionReport()) || {}
    tasks.value = toList(await getGeoGapTasks(50))
  } finally { loading.value = false }
}

const onDone = async (row) => {
  await markGeoTaskDone(row.id)
  ElMessage.success(t('auto.DecisionReport.fad5222c'))
  await load()
}

onMounted(load)
</script>
