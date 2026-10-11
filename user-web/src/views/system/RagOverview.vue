<template>
  <div class="rag-overview-container">
    <el-card class="rag-overview-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.RagOverview.d3406221') }}</span>
        </div>
      </template>
      
      <div class="overview-content">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-card class="stat-card">
              <div class="stat-item">
                <div class="stat-number">{{ stats.totalProducts }}</div>
                <div class="stat-label">{{ $t('auto.RagOverview.b0291c22') }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card class="stat-card">
              <div class="stat-item">
                <div class="stat-number">{{ stats.activeProducts }}</div>
                <div class="stat-label">{{ $t('auto.RagOverview.51a08641') }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card class="stat-card">
              <div class="stat-item">
                <div class="stat-number">{{ stats.activatedAccounts }}</div>
                <div class="stat-label">{{ $t('auto.RagOverview.31ceae5b') }}</div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <el-divider />

        <div class="quick-actions">
          <h3>{{ $t('auto.RagOverview.002df904') }}</h3>
          <el-row :gutter="20">
            <el-col :span="6">
              <el-card class="action-card" @click="goToConfig">
                <div class="action-item">
                  <el-icon class="action-icon"><Setting /></el-icon>
                  <div class="action-label">{{ $t('auto.RagOverview.924c17cf') }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="action-card" @click="goToAccountConfig">
                <div class="action-item">
                  <el-icon class="action-icon"><User /></el-icon>
                  <div class="action-label">{{ $t('auto.RagOverview.a16a35a5') }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="action-card" @click="goToAccountConfig">
                <div class="action-item">
                  <el-icon class="action-icon"><ChatDotRound /></el-icon>
                  <div class="action-label">{{ $t('auto.RagOverview.1f6097c4') }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="action-card" @click="goToDocs">
                <div class="action-item">
                  <el-icon class="action-icon"><Document /></el-icon>
                  <div class="action-label">{{ $t('auto.RagOverview.c39ced2c') }}</div>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <el-divider />

        <div class="recent-activity">
          <h3>{{ $t('auto.RagOverview.73378a1f') }}</h3>
          <el-table :data="activities" style="width: 100%">
            <el-table-column prop="time" :label="$t('auto.RagOverview.19fcb9eb')" width="180" />
            <el-table-column prop="action" :label="$t('auto.RagOverview.2b6bc0f2')" />
            <el-table-column prop="user" :label="$t('auto.RagOverview.1fd02a90')" width="120" />
            <el-table-column prop="status" :label="$t('auto.RagOverview.3fea7ca7')" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusTagType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </el-card>

    
    <el-dialog v-model="helpDialogVisible" :title="$t('auto.RagOverview.d6c1434b')" width="780px" :close-on-click-modal="false">
      <div class="help-doc-content">
        <el-collapse v-model="activeHelpSections">
          <el-collapse-item :title="$t('auto.RagOverview.40d5f4b2')" name="1">
            <div class="help-section">
              <p>{{ $t('auto.RagOverview.853b1a14') }}</p>
              <p><strong>{{ $t('auto.RagOverview.a3d37957') }}</strong></p>
              <ul>
                <li>{{ $t('auto.RagOverview.98e74b3b') }}</li>
                <li>{{ $t('auto.RagOverview.e39f91f3') }}</li>
                <li>{{ $t('auto.RagOverview.f78301c0') }}</li>
                <li>{{ $t('auto.RagOverview.081fb90f') }}</li>
              </ul>
            </div>
          </el-collapse-item>

          <el-collapse-item :title="$t('auto.RagOverview.321e9b7f')" name="2">
            <div class="help-section">
              <p><strong>{{ $t('auto.RagOverview.1cdf6381') }}</strong>{{ $t('auto.RagOverview.3769adad') }}</p>
              <p><strong>{{ $t('auto.RagOverview.e8ad9fab') }}</strong>{{ $t('auto.RagOverview.4531fbc8') }}</p>
              <ul>
                <li><strong>{{ $t('auto.RagOverview.0fcac790') }}</strong>{{ $t('auto.RagOverview.5496af86') }}</li>
                <li><strong>{{ $t('auto.RagOverview.b421bd00') }}</strong>{{ $t('auto.RagOverview.ac5641ca') }}</li>
                <li><strong>{{ $t('auto.RagOverview.29e965a4') }}</strong>{{ $t('auto.RagOverview.9b6fdbbb') }}</li>
                <li><strong>{{ $t('auto.RagOverview.97af7e41') }}</strong>{{ $t('auto.RagOverview.1832818f') }}</li>
                <li><strong>{{ $t('auto.RagOverview.577bc8c8') }}</strong>{{ $t('auto.RagOverview.9287f915') }}</li>
              </ul>
              <p><strong>{{ $t('auto.RagOverview.a5da0dc5') }}</strong>{{ $t('auto.RagOverview.70329e4c') }}</p>
            </div>
          </el-collapse-item>

          <el-collapse-item :title="$t('auto.RagOverview.ff7b0905')" name="3">
            <div class="help-section">
              <p><strong>{{ $t('auto.RagOverview.1cdf6381') }}</strong>{{ $t('auto.RagOverview.ad5afa7a') }}</p>
              <p><strong>{{ $t('auto.RagOverview.e8ad9fab') }}</strong>{{ $t('auto.RagOverview.37e560ae') }}</p>
              <ul>
                <li>{{ $t('auto.RagOverview.11f65b04') }}</li>
                <li>{{ $t('auto.RagOverview.56227f1c') }}</li>
                <li>{{ $t('auto.RagOverview.6ca9087d') }}</li>
                <li>{{ $t('auto.RagOverview.6bc0b5b3') }}</li>
                <li>{{ $t('auto.RagOverview.537a9b56') }}</li>
              </ul>
              <p><strong>{{ $t('auto.RagOverview.a5da0dc5') }}</strong>{{ $t('auto.RagOverview.9d540f8b') }}</p>
              <p><strong>{{ $t('auto.RagOverview.4cdc8893') }}</strong>{{ $t('auto.RagOverview.58b1846f') }}</p>
              <p><strong>{{ $t('auto.RagOverview.f35cf621') }}</strong></p>
              <ul>
                <li>{{ $t('auto.RagOverview.819c2b45') }}</li>
                <li>{{ $t('auto.RagOverview.0712c103') }}</li>
                <li>{{ $t('auto.RagOverview.6bba16db') }}</li>
                <li>{{ $t('auto.RagOverview.4e3a949b') }}</li>
              </ul>
            </div>
          </el-collapse-item>

          <el-collapse-item :title="$t('auto.RagOverview.98850364')" name="4">
            <div class="help-section">
              <p><strong>{{ $t('auto.RagOverview.1cdf6381') }}</strong>{{ $t('auto.RagOverview.e2d2b4b2') }}</p>
              <p><strong>{{ $t('auto.RagOverview.e8ad9fab') }}</strong>{{ $t('auto.RagOverview.4db3c80a') }}</p>
              <p><strong>{{ $t('auto.RagOverview.a5da0dc5') }}</strong>{{ $t('auto.RagOverview.ab11f796') }}</p>
              <ul>
                <li><strong>{{ $t('auto.RagOverview.68cf1f5e') }}</strong>{{ $t('auto.RagOverview.6759c68e') }}</li>
                <li><strong>{{ $t('auto.RagOverview.ef75959e') }}</strong>{{ $t('auto.RagOverview.58b6bff0') }}</li>
                <li><strong>{{ $t('auto.RagOverview.873fe7c1') }}</strong>{{ $t('auto.RagOverview.c6100798') }}</li>
                <li><strong>{{ $t('auto.RagOverview.aff5442d') }}</strong>{{ $t('auto.RagOverview.18a3e3a2') }}</li>
              </ul>
              <p><strong>{{ $t('auto.RagOverview.4cdc8893') }}</strong>{{ $t('auto.RagOverview.e5e8a89b') }}</p>
            </div>
          </el-collapse-item>

          <el-collapse-item :title="$t('auto.RagOverview.90ebc26e')" name="5">
            <div class="help-section">
              <p><strong>{{ $t('auto.RagOverview.4cc9efe8') }}</strong></p>
              <p>{{ $t('auto.RagOverview.bbf19b0c') }}</p>

              <p><strong>{{ $t('auto.RagOverview.25fdf50b') }}</strong></p>
              <p>{{ $t('auto.RagOverview.a20efa94') }}</p>

              <p><strong>{{ $t('auto.RagOverview.4114d088') }}</strong></p>
              <p>{{ $t('auto.RagOverview.8a2dee01') }}</p>

              <p><strong>{{ $t('auto.RagOverview.126a6a63') }}</strong></p>
              <p>{{ $t('auto.RagOverview.d310882a') }}</p>

              <p><strong>{{ $t('auto.RagOverview.024e92a1') }}</strong></p>
              <p>{{ $t('auto.RagOverview.fffe41d1') }}</p>
            </div>
          </el-collapse-item>
        </el-collapse>
      </div>
      <template #footer>
        <el-button type="primary" @click="helpDialogVisible = false">{{ $t('auto.RagOverview.b15d9127') }}</el-button>
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
const STATUS_ZH_LABEL = { '成功': t('auto.RagOverview.330363df'), '失败': t('auto.RagOverview.acd5cb84'), '运行中': t('auto.RagOverview.fb852fc6') }
const getStatusLabel = (status) => STATUS_ZH_LABEL[status] || (status === 'success' ? t('auto.RagOverview.330363df') : status === 'failed' ? t('auto.RagOverview.acd5cb84') : status === 'running' ? t('auto.RagOverview.fb852fc6') : status)

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
      console.warn(t('auto.RagOverview.3e914998'), e);
    }

    try {
      const accountsRes = await ragProductConfigAPI.getAccountConfig({})
      const accounts = accountsRes?.list || accountsRes || []
      if (Array.isArray(accounts)) {
        stats.value.activatedAccounts = accounts.filter((a) => a.is_active !== false).length
      }
    } catch (e) {
      console.warn(t('auto.RagOverview.139fa04d'), e);
    }

    if (Array.isArray(products) && products.length > 0) {
      activities.value = products.slice(0, 5).map((p) => ({
        time: p.updated_at || p.created_at || t('auto.RagOverview.9418d7cb'),
        action: `RAG产品: ${p.name || t('auto.RagOverview.7f0425a8')}`,
        user: p.creator || 'system',
        status: p.is_active ? t('auto.RagOverview.330363df') : t('auto.RagOverview.acd5cb84')
      }))
    } else {
      activities.value = []
    }
  } catch (error) {
    console.error(t('auto.RagOverview.b6ee622f'), error)
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