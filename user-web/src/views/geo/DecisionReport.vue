<template>
  <div class="geo-page">

    <div class="p-4">
    <el-row :gutter="16" class="mb-4">
      <el-col :span="8"><el-card><el-statistic :title="$t('auto.DecisionReport.1.6ff099')" :value="report.leads_captured || 0" /></el-card></el-col>
      <el-col :span="8"><el-card><el-statistic :title="$t('auto.DecisionReport.2.460225')" :value="report.tasks_pending || 0" /></el-card></el-col>
      <el-col :span="8"><el-card><el-statistic :title="$t('auto.DecisionReport.3.08fd57')" :value="report.tasks_done || 0" /></el-card></el-col>
    </el-row>

    <el-card>
      <template #header><span>{{ $t('auto.DecisionReport.9.397bda') }}</span></template>
      <el-table :data="tasks" v-loading="loading" size="small">
        <el-table-column prop="keyword" :label="$t('auto.DecisionReport.4.9699a5')" min-width="140" />
        <el-table-column prop="intent" :label="$t('auto.DecisionReport.5.764aca')" width="90" />
        <el-table-column prop="gap_type" :label="$t('auto.DecisionReport.6.e755cd')" width="150" />
        <el-table-column prop="detail" :label="$t('auto.DecisionReport.7.f26225')" min-width="260" show-overflow-tooltip />
        <el-table-column :label="$t('auto.DecisionReport.8.2b6bc0')" width="100">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="onDone(row)">{{ $t('auto.DecisionReport.10.769d88') }}</el-button>
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
  ElMessage.success(t('auto.DecisionReport.11.fad522'))
  await load()
}

onMounted(load)
</script>
