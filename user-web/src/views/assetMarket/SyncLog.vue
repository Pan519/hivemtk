<template>
  <div class="sync-log">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.SyncLog.8.e1be44') }}</span>
          <el-button @click="$router.back()">{{ $t('auto.SyncLog.9.5f4112') }}</el-button>
        </div>
      </template>
      <el-form inline>
        <el-form-item :label="$t('auto.SyncLog.1.69847d')">
          <el-input v-model="assetId" clearable :placeholder="$t('auto.SyncLog.2.c20cba')" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">{{ $t('auto.SyncLog.10.bee912') }}</el-button>
        </el-form-item>
      </el-form>
      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="created_at" :label="$t('auto.SyncLog.3.19fcb9')" width="180" />
        <el-table-column prop="asset_id" :label="$t('auto.SyncLog.4.5110a0')" min-width="180" />
        <el-table-column prop="action" :label="$t('auto.SyncLog.5.2b6bc0')" width="120" />
        <el-table-column prop="status" :label="$t('auto.SyncLog.6.3fea7c')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="error_msg" :label="$t('auto.SyncLog.7.4604d5')" min-width="200" show-overflow-tooltip />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { syncLog } from '@/api/assetMarket'

const { t } = useI18n()

const loading = ref(false)
const list = ref([])
const assetId = ref('')

const fetchList = async () => {
  loading.value = true
  try {
    const resp = await syncLog({ asset_id: assetId.value, limit: 100 })
    list.value = resp?.data || resp || []
  } finally {
    loading.value = false
  }
}

onMounted(fetchList)
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
