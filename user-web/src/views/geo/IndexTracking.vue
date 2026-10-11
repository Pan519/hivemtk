<template>
  <div class="index-tracking">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.IndexTracking.e4ed9651') }}</span>
          <div>
            <el-button type="primary" @click="load">{{ $t('auto.IndexTracking.694fc5ef') }}</el-button>
            <el-button @click="verifyAll">{{ $t('auto.IndexTracking.ffb91a32') }}</el-button>
          </div>
        </div>
      </template>
      <el-empty v-if="!trackings.length" :description="$t('auto.IndexTracking.9286461a')" />
      <el-table v-else :data="trackings" size="small" border>
        <el-table-column prop="url" label="URL" show-overflow-tooltip />
        <el-table-column prop="engine" :label="$t('auto.IndexTracking.6c1b19dd')" width="100" />
        <el-table-column :label="$t('auto.IndexTracking.1ed3c154')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.indexed ? 'success' : 'info'" size="small">{{ row.indexed ? '✅' : '—' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="rank_position" :label="$t('auto.IndexTracking.a4dc00d8')" width="80" />
        <el-table-column :label="$t('auto.IndexTracking.34647732')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.ai_cited ? 'success' : 'warning'" size="small">
              {{ row.ai_cited ? '✅ ' + row.ai_cite_count + '次' : '❌' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_checked" :label="$t('auto.IndexTracking.c1ce231c')" width="160" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { geoApi } from '@/api/geo'
import { ElMessage } from 'element-plus'
const { t } = useI18n()

const trackings = ref([])

async function load() {
  try {
    const res = await geoApi.getIndexFunnel()
    trackings.value = res?.trackings || []
  } catch (e) { /* 忽略异常：失败时保持既有状态，不打断用户 */ }
}
async function verifyAll() {
  try {
    await geoApi.verifyAllIndex()
    ElMessage.success(t('auto.IndexTracking.f99996de'))
    load()
  } catch (e) { ElMessage.error(t('auto.IndexTracking.acd5cb84')) }
}
onMounted(load)
</script>
