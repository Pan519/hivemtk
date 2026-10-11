<template>
  <div class="tool-management-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ $t('auto.ToolList.31.d3c67f') }}</h2>
          <p class="subtitle">{{ $t('auto.ToolList.32.949f0e') }}</p>
        </div>
        <div class="header-actions">
          <el-button @click="loadTools" :loading="loading">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    
    <el-row :gutter="16" class="stats-row">
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-value">{{ stats.total }}</div>
          <div class="stat-label">{{ $t('auto.ToolList.33.be7fc8') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat-card stat-success">
          <div class="stat-value">{{ stats.enabled }}</div>
          <div class="stat-label">{{ $t('auto.ToolList.10.53ace4') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat-card stat-warning">
          <div class="stat-value">{{ stats.disabled }}</div>
          <div class="stat-label">{{ $t('auto.ToolList.11.1c1ed9') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat-card stat-info">
          <div class="stat-value">{{ stats.categories }}</div>
          <div class="stat-label">{{ $t('auto.ToolList.34.b0a53c') }}</div>
        </el-card>
      </el-col>
    </el-row>

    
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filter">
        <el-form-item :label="$t('auto.ToolList.1.d0771a')">
          <el-select v-model="filter.category" :placeholder="$t('auto.ToolList.2.1a7503')" clearable @change="loadTools">
            <el-option :label="$t('auto.ToolList.3.59fa82')" value="reach" />
            <el-option :label="$t('auto.ToolList.4.110f30')" value="customer" />
            <el-option :label="$t('auto.ToolList.5.b9a1ca')" value="knowledge" />
            <el-option :label="$t('auto.ToolList.6.7a8e9c')" value="business" />
            <el-option :label="$t('auto.ToolList.7.358d33')" value="private_message" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.ToolList.8.3fea7c')">
          <el-select v-model="filter.enabled" :placeholder="$t('auto.ToolList.9.443483')" clearable @change="loadTools">
            <el-option :label="$t('auto.ToolList.10.53ace4')" :value="true" />
            <el-option :label="$t('auto.ToolList.11.1c1ed9')" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadTools">{{ $t('auto.ToolList.35.e5f71f') }}</el-button>
          <el-button @click="resetFilter">{{ $t('auto.ToolList.36.4b9c32') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    
    <el-card shadow="never">
      <el-table :data="tools" v-loading="loading" stripe border>
        <el-table-column :label="$t('auto.ToolList.12.20dce2')" min-width="280">
          <template #default="{ row }">
            <div class="tool-info">
              <div class="tool-name">{{ row.tool_name }}</div>
              <div class="tool-desc">{{ row.config?.description_zh || getCategoryLabel(row.category) }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.ToolList.1.d0771a')" width="120">
          <template #default="{ row }">
            <el-tag :type="getCategoryTagType(row.category)" size="small">
              {{ getCategoryLabel(row.category) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.ToolList.8.3fea7c')" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_enabled"
              @change="onToggleTool(row)"
              :loading="row._toggling"
            />
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.ToolList.13.0529dc')" min-width="200">
          <template #default="{ row }">
            <div v-if="row.bound_accounts?.length">
              <el-tag
                v-for="(account, index) in row.bound_accounts"
                :key="index"
                type="success"
                size="small"
                class="account-tag"
              >
                {{ account.account_id }}
                <el-icon v-if="account.is_primary" class="primary-icon"><Star /></el-icon>
              </el-tag>
            </div>
            <div v-else-if="needsAccount(row.category)">
              <el-button link type="primary" size="small" @click="openBindDialog(row)">
                绑定账号
              </el-button>
            </div>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.ToolList.14.2b6bc0')" width="120" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="openDetail(row)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <el-dialog v-model="detailVisible" :title="$t('auto.ToolList.15.ab7eba')" width="600px">
      <el-descriptions :column="2" border v-if="currentTool">
        <el-descriptions-item :label="$t('auto.ToolList.16.ac6116')">{{ currentTool.tool_name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.ToolList.1.d0771a')">
          <el-tag :type="getCategoryTagType(currentTool.category)" size="small">
            {{ getCategoryLabel(currentTool.category) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.ToolList.8.3fea7c')">
          <el-tag :type="currentTool.is_enabled ? 'success' : 'info'" size="small">
            {{ currentTool.is_enabled ? '已启用' : '已禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.ToolList.17.0a05b1')">{{ currentTool.display_order }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.ToolList.18.224e2c')" :span="2">
          <pre class="config-json">{{ JSON.stringify(currentTool.config || {}, null, 2) }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    
    <el-dialog v-model="bindDialogVisible" :title="`绑定账号 - ${currentTool?.tool_name}`" width="600px">
      <el-form :model="bindForm" label-width="100px">
        <el-form-item :label="$t('auto.ToolList.19.c65d92')">
          <el-select v-model="bindForm.account_type" :placeholder="$t('auto.ToolList.20.c3b31c')">
            <el-option :label="$t('auto.ToolList.21.485c3a')" value="sms" />
            <el-option :label="$t('auto.ToolList.22.e9e805')" value="email" />
            <el-option label="Telegram" value="telegram" />
            <el-option :label="$t('auto.ToolList.23.ff17b9')" value="wecom" />
            <el-option :label="$t('auto.ToolList.24.7714e5')" value="feishu" />
            <el-option label="WhatsApp" value="whatsapp" />
            <el-option :label="$t('auto.ToolList.25.277a98')" value="douyin" />
            <el-option :label="$t('auto.ToolList.26.f32c0a')" value="kuaishou" />
            <el-option :label="$t('auto.ToolList.27.ae26bb')" value="xhs" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.ToolList.28.14d171')">
          <el-input v-model="bindForm.account_id" :placeholder="$t('auto.ToolList.29.37c32e')" />
        </el-form-item>
        <el-form-item :label="$t('auto.ToolList.30.a2c1a2')">
          <el-switch v-model="bindForm.is_primary" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="bindDialogVisible = false">{{ $t('auto.ToolList.37.625fb2') }}</el-button>
        <el-button type="primary" @click="onBind">{{ $t('auto.ToolList.38.38cf16') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Refresh, Star } from '@element-plus/icons-vue'
import { listTools, updateToolStatus, bindToolAccount } from '@/api/aiTool.js'

const { t } = useI18n()

const loading = ref(false);
const tools = ref([])
const filter = reactive({
  category: '',
  enabled: null
})
const stats = reactive({
  total: 0,
  enabled: 0,
  disabled: 0,
  categories: 5
})

const detailVisible = ref(false);
const currentTool = ref(null)

const bindDialogVisible = ref(false);
const bindForm = reactive({
  account_type: '',
  account_id: '',
  is_primary: false
})

const categoryMap = {
  reach: t('auto.ToolList.3.59fa82'),
  customer: t('auto.ToolList.4.110f30'),
  knowledge: t('auto.ToolList.5.b9a1ca'),
  business: t('auto.ToolList.6.7a8e9c'),
  private_message: t('auto.ToolList.7.358d33')
};

const getCategoryLabel = (category) => categoryMap[category] || category

const getCategoryTagType = (category) => {
  const map = {
    reach: 'primary',
    customer: 'success',
    knowledge: 'warning',
    business: 'info',
    private_message: 'danger'
  }
  return map[category] || 'info'
}

const needsAccount = (category) => {
  return category === 'reach' || category === 'private_message'
};

const loadTools = async () => {
  loading.value = true
  try {
    const params = {}
    if (filter.category) params.category = filter.category
    if (filter.enabled !== null) params.enabled = filter.enabled

    const res = await listTools(params)
    tools.value = res?.list || []

    stats.total = tools.value.length;
    stats.enabled = tools.value.filter(t => t.is_enabled).length
    stats.disabled = tools.value.filter(t => !t.is_enabled).length
  } catch (e) {
    ElMessage.error(t('auto.ToolList.39.162848') + e.message)
  } finally {
    loading.value = false
  }
};

const onToggleTool = async (tool) => {
  tool._toggling = true
  try {
    await updateToolStatus(tool.tool_name, tool.is_enabled)
    ElMessage.success(`${tool.is_enabled ? t('auto.ToolList.40.7854b5') : t('auto.ToolList.41.710ad0')}成功`)
  } catch (e) {
    tool.is_enabled = !tool.is_enabled
    ElMessage.error(t('auto.ToolList.42.79585f') + e.message)
  } finally {
    tool._toggling = false
  }
};

const resetFilter = () => {
  filter.category = ''
  filter.enabled = null
  loadTools()
};

const openDetail = (tool) => {
  currentTool.value = tool
  detailVisible.value = true
};

const openBindDialog = (tool) => {
  currentTool.value = tool
  bindForm.account_type = ''
  bindForm.account_id = ''
  bindForm.is_primary = false
  bindDialogVisible.value = true
};

const onBind = async () => {
  if (!bindForm.account_type || !bindForm.account_id) {
    ElMessage.warning(t('auto.ToolList.43.d2d422'))
    return
  }

  try {
    await bindToolAccount(
      currentTool.value.tool_name,
      bindForm.account_type,
      bindForm.account_id,
      bindForm.is_primary
    )
    ElMessage.success(t('auto.ToolList.44.1974fe'))
    bindDialogVisible.value = false
    loadTools()
  } catch (e) {
    ElMessage.error(t('auto.ToolList.45.8cbf8b') + e.message)
  }
};

onMounted(() => {
  loadTools()
});
</script>

<style scoped lang="scss">
.tool-management-page {
  padding: 20px;
}

.header-card {
  margin-bottom: 16px;
  
  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    
    h2 {
      margin: 0 0 6px 0;
      font-size: 20px;
    }
    
    .subtitle {
      color: #909399;
      margin: 0;
      font-size: 13px;
    }
  }
}

.stats-row {
  margin-bottom: 16px;
}

.stat-card {
  text-align: center;
  
  .stat-value {
    font-size: 28px;
    font-weight: bold;
    color: #303133;
  }
  
  .stat-label {
    font-size: 13px;
    color: #909399;
    margin-top: 4px;
  }
  
  &.stat-success .stat-value { color: #67C23A; }
  &.stat-warning .stat-value { color: #E6A23C; }
  &.stat-info .stat-value { color: #409EFF; }
}

.filter-card {
  margin-bottom: 16px;
}

.tool-info {
  .tool-name {
    font-weight: 500;
    color: #303133;
  }
  
  .tool-desc {
    font-size: 12px;
    color: #909399;
    margin-top: 4px;
  }
}

.account-tag {
  margin-right: 4px;
  margin-bottom: 4px;
  
  .primary-icon {
    margin-left: 2px;
    color: #F7BA2A;
  }
}

.text-muted {
  color: #C0C4CC;
}

.config-json {
  background: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  font-size: 12px;
  max-height: 200px;
  overflow-y: auto;
}
</style>