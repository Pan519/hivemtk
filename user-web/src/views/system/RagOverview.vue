<template>
  <div class="rag-overview-container">
    <el-card class="rag-overview-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.RagOverview.11.d34062') }}</span>
        </div>
      </template>
      
      <div class="overview-content">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-card class="stat-card">
              <div class="stat-item">
                <div class="stat-number">{{ stats.totalProducts }}</div>
                <div class="stat-label">{{ $t('auto.RagOverview.12.b0291c') }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card class="stat-card">
              <div class="stat-item">
                <div class="stat-number">{{ stats.activeProducts }}</div>
                <div class="stat-label">{{ $t('auto.RagOverview.13.51a086') }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card class="stat-card">
              <div class="stat-item">
                <div class="stat-number">{{ stats.activatedAccounts }}</div>
                <div class="stat-label">{{ $t('auto.RagOverview.14.31ceae') }}</div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <el-divider />

        <div class="quick-actions">
          <h3>{{ $t('auto.RagOverview.15.002df9') }}</h3>
          <el-row :gutter="20">
            <el-col :span="6">
              <el-card class="action-card" @click="goToConfig">
                <div class="action-item">
                  <el-icon class="action-icon"><Setting /></el-icon>
                  <div class="action-label">{{ $t('auto.RagOverview.16.924c17') }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="action-card" @click="goToAccountConfig">
                <div class="action-item">
                  <el-icon class="action-icon"><User /></el-icon>
                  <div class="action-label">{{ $t('auto.RagOverview.17.a16a35') }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="action-card" @click="goToAccountConfig">
                <div class="action-item">
                  <el-icon class="action-icon"><ChatDotRound /></el-icon>
                  <div class="action-label">{{ $t('auto.RagOverview.18.1f6097') }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="action-card" @click="goToDocs">
                <div class="action-item">
                  <el-icon class="action-icon"><Document /></el-icon>
                  <div class="action-label">{{ $t('auto.RagOverview.19.c39ced') }}</div>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <el-divider />

        <div class="recent-activity">
          <h3>{{ $t('auto.RagOverview.20.73378a') }}</h3>
          <el-table :data="activities" style="width: 100%">
            <el-table-column prop="time" :label="$t('auto.RagOverview.1.19fcb9')" width="180" />
            <el-table-column prop="action" :label="$t('auto.RagOverview.2.2b6bc0')" />
            <el-table-column prop="user" :label="$t('auto.RagOverview.3.1fd02a')" width="120" />
            <el-table-column prop="status" :label="$t('auto.RagOverview.4.3fea7c')" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusTagType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </el-card>

    
    <el-dialog v-model="helpDialogVisible" :title="$t('auto.RagOverview.5.d6c143')" width="780px" :close-on-click-modal="false">
      <div class="help-doc-content">
        <el-collapse v-model="activeHelpSections">
          <el-collapse-item :title="$t('auto.RagOverview.6.40d5f4')" name="1">
            <div class="help-section">
              <p>{{ $t('auto.RagOverview.21.853b1a') }}</p>
              <p><strong>{{ $t('auto.RagOverview.22.a3d379') }}</strong></p>
              <ul>
                <li>{{ $t('auto.RagOverview.23.98e74b') }}</li>
                <li>{{ $t('auto.RagOverview.24.e39f91') }}</li>
                <li>{{ $t('auto.RagOverview.25.f78301') }}</li>
                <li>{{ $t('auto.RagOverview.26.081fb9') }}</li>
              </ul>
            </div>
          </el-collapse-item>

          <el-collapse-item :title="$t('auto.RagOverview.7.321e9b')" name="2">
            <div class="help-section">
              <p><strong>{{ $t('auto.RagOverview.27.1cdf63') }}</strong>{{ $t('auto.RagOverview.28.3769ad') }}</p>
              <p><strong>{{ $t('auto.RagOverview.29.e8ad9f') }}</strong>{{ $t('auto.RagOverview.30.4531fb') }}</p>
              <ul>
                <li><strong>{{ $t('auto.RagOverview.31.0fcac7') }}</strong>{{ $t('auto.RagOverview.32.5496af') }}</li>
                <li><strong>{{ $t('auto.RagOverview.33.b421bd') }}</strong>{{ $t('auto.RagOverview.34.ac5641') }}</li>
                <li><strong>{{ $t('auto.RagOverview.35.29e965') }}</strong>{{ $t('auto.RagOverview.36.9b6fdb') }}</li>
                <li><strong>{{ $t('auto.RagOverview.37.97af7e') }}</strong>{{ $t('auto.RagOverview.38.183281') }}</li>
                <li><strong>{{ $t('auto.RagOverview.39.577bc8') }}</strong>{{ $t('auto.RagOverview.40.9287f9') }}</li>
              </ul>
              <p><strong>{{ $t('auto.RagOverview.41.a5da0d') }}</strong>{{ $t('auto.RagOverview.42.70329e') }}</p>
            </div>
          </el-collapse-item>

          <el-collapse-item :title="$t('auto.RagOverview.8.ff7b09')" name="3">
            <div class="help-section">
              <p><strong>{{ $t('auto.RagOverview.27.1cdf63') }}</strong>{{ $t('auto.RagOverview.43.ad5afa') }}</p>
              <p><strong>{{ $t('auto.RagOverview.29.e8ad9f') }}</strong>{{ $t('auto.RagOverview.44.37e560') }}</p>
              <ul>
                <li>{{ $t('auto.RagOverview.45.11f65b') }}</li>
                <li>{{ $t('auto.RagOverview.46.56227f') }}</li>
                <li>{{ $t('auto.RagOverview.47.6ca908') }}</li>
                <li>{{ $t('auto.RagOverview.48.6bc0b5') }}</li>
                <li>{{ $t('auto.RagOverview.49.537a9b') }}</li>
              </ul>
              <p><strong>{{ $t('auto.RagOverview.41.a5da0d') }}</strong>{{ $t('auto.RagOverview.50.9d540f') }}</p>
              <p><strong>{{ $t('auto.RagOverview.51.4cdc88') }}</strong>{{ $t('auto.RagOverview.52.58b184') }}</p>
              <p><strong>{{ $t('auto.RagOverview.53.f35cf6') }}</strong></p>
              <ul>
                <li>{{ $t('auto.RagOverview.54.819c2b') }}</li>
                <li>{{ $t('auto.RagOverview.55.0712c1') }}</li>
                <li>{{ $t('auto.RagOverview.56.6bba16') }}</li>
                <li>{{ $t('auto.RagOverview.57.4e3a94') }}</li>
              </ul>
            </div>
          </el-collapse-item>

          <el-collapse-item :title="$t('auto.RagOverview.9.988503')" name="4">
            <div class="help-section">
              <p><strong>{{ $t('auto.RagOverview.27.1cdf63') }}</strong>{{ $t('auto.RagOverview.58.e2d2b4') }}</p>
              <p><strong>{{ $t('auto.RagOverview.29.e8ad9f') }}</strong>{{ $t('auto.RagOverview.59.4db3c8') }}</p>
              <p><strong>{{ $t('auto.RagOverview.41.a5da0d') }}</strong>{{ $t('auto.RagOverview.60.ab11f7') }}</p>
              <ul>
                <li><strong>{{ $t('auto.RagOverview.61.68cf1f') }}</strong>{{ $t('auto.RagOverview.62.6759c6') }}</li>
                <li><strong>{{ $t('auto.RagOverview.63.ef7595') }}</strong>{{ $t('auto.RagOverview.64.58b6bf') }}</li>
                <li><strong>{{ $t('auto.RagOverview.65.873fe7') }}</strong>{{ $t('auto.RagOverview.66.c61007') }}</li>
                <li><strong>{{ $t('auto.RagOverview.67.aff544') }}</strong>{{ $t('auto.RagOverview.68.18a3e3') }}</li>
              </ul>
              <p><strong>{{ $t('auto.RagOverview.51.4cdc88') }}</strong>{{ $t('auto.RagOverview.69.e5e8a8') }}</p>
            </div>
          </el-collapse-item>

          <el-collapse-item :title="$t('auto.RagOverview.10.90ebc2')" name="5">
            <div class="help-section">
              <p><strong>{{ $t('auto.RagOverview.70.4cc9ef') }}</strong></p>
              <p>{{ $t('auto.RagOverview.71.bbf19b') }}</p>

              <p><strong>{{ $t('auto.RagOverview.72.25fdf5') }}</strong></p>
              <p>{{ $t('auto.RagOverview.73.a20efa') }}</p>

              <p><strong>{{ $t('auto.RagOverview.74.4114d0') }}</strong></p>
              <p>{{ $t('auto.RagOverview.75.8a2dee') }}</p>

              <p><strong>{{ $t('auto.RagOverview.76.126a6a') }}</strong></p>
              <p>{{ $t('auto.RagOverview.77.d31088') }}</p>

              <p><strong>{{ $t('auto.RagOverview.78.024e92') }}</strong></p>
              <p>{{ $t('auto.RagOverview.79.fffe41') }}</p>
            </div>
          </el-collapse-item>
        </el-collapse>
      </div>
      <template #footer>
        <el-button type="primary" @click="helpDialogVisible = false">{{ $t('auto.RagOverview.80.b15d91') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ragProductConfigAPI } from '@/api/ragProductConfig'
import { knowledgeBaseAPI } from '@/api/knowledgeBase'

const router = useRouter()

const { t } = useI18n()

const stats = ref({
  totalProducts: 0,
  activeProducts: 0,
  activatedAccounts: 0,
  totalDocuments: 0
});

const activities = ref([]);

const STATUS_ZH_TO_VALUE = { '成功': 'success', '失败': 'failed', '运行中': 'running' };
const getStatusTagType = (status) => {
  const v = STATUS_ZH_TO_VALUE[status] || status
  switch (v) {
    case 'success': return 'success'
    case 'failed':  return 'danger'
    case 'running': return 'warning'
    default:        return 'info'
  }
}
const STATUS_ZH_LABEL = { '成功': t('auto.RagOverview.81.330363'), '失败': t('auto.RagOverview.82.acd5cb'), '运行中': t('auto.RagOverview.83.fb852f') }
const getStatusLabel = (status) => STATUS_ZH_LABEL[status] || (status === 'success' ? t('auto.RagOverview.81.330363') : status === 'failed' ? t('auto.RagOverview.82.acd5cb') : status === 'running' ? t('auto.RagOverview.83.fb852f') : status)

const goToConfig = () => {
  router.push('/system/rag-product-config')
};

const goToAccountConfig = () => {
  router.push('/system/rag-product-config')
};

const helpDialogVisible = ref(false);
const activeHelpSections = ref(['1'])

const goToDocs = () => {
  helpDialogVisible.value = true
};

const loadStats = async () => {
  try {
    const productsRes = await ragProductConfigAPI.getRagProducts({});
    const products = productsRes?.items || productsRes?.list || []
    if (Array.isArray(products)) {
      stats.value.totalProducts = products.length
      stats.value.activeProducts = products.filter((p) => p.is_active !== false).length
    }

    try {
      const kbsRes = await knowledgeBaseAPI.listKBs({ page: 1, page_size: 100 })
      const kbs = kbsRes?.list || kbsRes?.items || (Array.isArray(kbsRes) ? kbsRes : [])
      stats.value.totalDocuments = kbs.reduce((sum, k) => sum + Number(k?.doc_count || k?.document_count || 0), 0)
    } catch (e) {
      console.warn(t('auto.RagOverview.84.3e9149'), e);
    }

    try {
      const accountsRes = await ragProductConfigAPI.getAccountConfig({})
      const accounts = accountsRes?.list || accountsRes || []
      if (Array.isArray(accounts)) {
        stats.value.activatedAccounts = accounts.filter((a) => a.is_active !== false).length
      }
    } catch (e) {
      console.warn(t('auto.RagOverview.85.139fa0'), e);
    }

    if (Array.isArray(products) && products.length > 0) {
      activities.value = products.slice(0, 5).map((p) => ({
        time: p.updated_at || p.created_at || t('auto.RagOverview.86.9418d7'),
        action: `RAG产品: ${p.name || t('auto.RagOverview.87.7f0425')}`,
        user: p.creator || 'system',
        status: p.is_active ? t('auto.RagOverview.81.330363') : t('auto.RagOverview.82.acd5cb')
      }))
    } else {
      activities.value = []
    }
  } catch (error) {
    console.error(t('auto.RagOverview.88.b6ee62'), error)
    ElMessage.warning(i18n.global.t('RAG 概览数据加载失败,请稍后重试'))
  }
};

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.rag-overview-container {
  padding: 20px;
}

.rag-overview-card {
  min-height: 600px;
}

.card-header {
  font-size: 18px;
  font-weight: bold;
}

.stat-item {
  text-align: center;
}

.stat-number {
  font-size: 28px;
  font-weight: bold;
  color: #4F46E5;
  margin-bottom: 8px;
}

.stat-label {
  color: #909399;
  font-size: 14px;
}

.stat-card {
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
}

.quick-actions h3 {
  margin-bottom: 20px;
}

.action-card {
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.action-icon {
  font-size: 32px;
  color: #4F46E5;
  margin-bottom: 8px;
}

.action-label {
  font-size: 14px;
  color: #606266;
}

.recent-activity h3 {
  margin-bottom: 20px;
}

.help-doc-content {
  max-height: 60vh;
  overflow-y: auto;
}

.help-section {
  padding: 10px 0;
  line-height: 1.8;
  color: #303133;
}

.help-section p {
  margin-bottom: 10px;
}

.help-section ul {
  margin-left: 20px;
  margin-bottom: 10px;
}

.help-section ul li {
  margin-bottom: 4px;
}
</style>