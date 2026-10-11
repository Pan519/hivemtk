<template>
  <div class="rag-product-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.RagProductManagement.99d89602') }}</span>
          <el-button 
            type="primary" 
            @click="showCreateDialog"
          >
            {{ $t('新增产品') }}
          </el-button>
        </div>
      </template>
      
      <el-table 
        :data="products" 
        style="width: 100%"
        v-loading="loading"
      >
        <el-table-column prop="name" :label="$t('auto.RagProductManagement.ff685527')" width="150"></el-table-column>
        <el-table-column prop="category" :label="$t('auto.RagProductManagement.383a6d16')" width="100">
          <template #default="{ row }">
            <el-tag :type="getCategoryTagType(row.category)">
              {{ getCategoryLabel(row.category) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="llm_model" :label="$t('auto.RagProductManagement.248d0fb8')" width="150"></el-table-column>
        <el-table-column prop="temperature" :label="$t('auto.RagProductManagement.723e5470')" width="100">
          <template #default="{ row }">
            {{ row.temperature.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="max_tokens" :label="$t('auto.RagProductManagement.b20b409f')" width="100">
          <template #default="{ row }">
            {{ row.max_tokens }}
          </template>
        </el-table-column>
        <el-table-column prop="top_p" label="Top-P" width="100">
          <template #default="{ row }">
            {{ row.top_p.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="description" :label="$t('auto.RagProductManagement.3bdd08ad')" show-overflow-tooltip></el-table-column>
        <el-table-column :label="$t('auto.RagProductManagement.3fea7ca7')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_active ? 'success' : 'info'">
              {{ row.is_active ? '激活' : '未激活' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.RagProductManagement.eca37cb0')" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.RagProductManagement.2b6bc0f2')" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="editProduct(row)">{{ $t('auto.RagProductManagement.95b351c8') }}</el-button>
            <el-button 
              size="small" 
              :type="row.is_active ? 'warning' : 'success'"
              @click="toggleActive(row)"
            >
              {{ row.is_active ? '停用' : '激活' }}
            </el-button>
            <el-button size="small" type="danger" @click="deleteProduct(row.id)">{{ $t('auto.RagProductManagement.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        style="margin-top: 20px; text-align: right;"
      />
    </el-card>
    
    
    <el-dialog 
      v-model="dialogVisible" 
      :title="dialogTitle" 
      width="500px"
      @close="resetForm"
    >
      <el-form 
        :model="form" 
        :rules="rules" 
        ref="formRef"
        label-width="100px"
      >
        <el-form-item :label="$t('auto.RagProductManagement.ff685527')" prop="name">
          <el-input 
            v-model="form.name" 
            :placeholder="$t('auto.RagProductManagement.448ab9f6')"
            :disabled="!!form.id"
          ></el-input>
        </el-form-item>
        
        <el-form-item :label="$t('auto.RagProductManagement.8d542597')" prop="category">
          <el-input 
            v-model="form.category" 
            :placeholder="$t('auto.RagProductManagement.0bae2a00')"
            style="width: 100%;"
          ></el-input>
        </el-form-item>
        
        <el-form-item :label="$t('auto.RagProductManagement.248d0fb8')">
          <el-input 
            v-model="form.llm_model" 
            :placeholder="$t('auto.RagProductManagement.4b287899')"
          ></el-input>
        </el-form-item>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('auto.RagProductManagement.723e5470')" prop="temperature">
              <el-slider
                v-model="form.temperature"
                :min="0"
                :max="2"
                :step="0.01"
                show-input
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.RagProductManagement.19c9a101')" prop="top_p">
              <el-slider
                v-model="form.top_p"
                :min="0"
                :max="1"
                :step="0.01"
                show-input
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('auto.RagProductManagement.1951437c')" prop="max_tokens">
              <el-input-number
                v-model="form.max_tokens"
                :min="1"
                :max="4096"
                style="width: 100%;"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.RagProductManagement.b8f0c714')" prop="frequency_penalty">
              <el-slider
                v-model="form.frequency_penalty"
                :min="-2"
                :max="2"
                :step="0.01"
                show-input
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item :label="$t('auto.RagProductManagement.ca5ae49a')" prop="presence_penalty">
          <el-slider
            v-model="form.presence_penalty"
            :min="-2"
            :max="2"
            :step="0.01"
            show-input
          />
        </el-form-item>
        
        <el-divider content-position="left">{{ $t('auto.RagProductManagement.bd540b21') }}</el-divider>

        
        <el-divider content-position="left">{{ $t('auto.RagProductManagement.13ada8a6') }}</el-divider>
        <el-form-item label="API BaseURL">
          <el-input v-model="form.llm_provider_config.base_url" :placeholder="$t('auto.RagProductManagement.dd75d37f')" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="form.llm_provider_config.api_key" type="password" show-password placeholder="sk-..." />
        </el-form-item>
        <el-form-item :label="$t('auto.RagProductManagement.8000f187')">
          <el-input v-model="form.llm_provider_config.model" :placeholder="$t('auto.RagProductManagement.b67062ae')" />
        </el-form-item>

        
        <el-divider content-position="left">{{ $t('auto.RagProductManagement.8301928d') }}</el-divider>
        <el-form-item label="API BaseURL">
          <el-input v-model="form.embedding_provider_config.base_url" :placeholder="$t('auto.RagProductManagement.7a18dd60')" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="form.embedding_provider_config.api_key" type="password" show-password placeholder="sk-..." />
        </el-form-item>
        <el-form-item :label="$t('auto.RagProductManagement.8000f187')">
          <el-input v-model="form.embedding_provider_config.model" :placeholder="$t('auto.RagProductManagement.8a31b178')" />
        </el-form-item>
        <el-form-item :label="$t('auto.RagProductManagement.f29c543f')">
          <el-input-number v-model="form.embedding_provider_config.dimension" :min="256" :max="4096" :step="128" />
        </el-form-item>
        <el-form-item :label="$t('auto.RagProductManagement.7854b52a')">
          <el-switch v-model="form.embedding_provider_config.enabled" />
        </el-form-item>

        
        <el-divider content-position="left">{{ $t('auto.RagProductManagement.980d3ebd') }}</el-divider>
        <el-form-item label="API BaseURL">
          <el-input v-model="form.rerank_provider_config.base_url" :placeholder="$t('auto.RagProductManagement.7a18dd60')" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="form.rerank_provider_config.api_key" type="password" show-password placeholder="sk-..." />
        </el-form-item>
        <el-form-item :label="$t('auto.RagProductManagement.8000f187')">
          <el-input v-model="form.rerank_provider_config.model" :placeholder="$t('auto.RagProductManagement.742fab31')" />
        </el-form-item>
        <el-form-item :label="$t('auto.RagProductManagement.7854b52a')">
          <el-switch v-model="form.rerank_provider_config.enabled" />
        </el-form-item>

        <el-form-item :label="$t('auto.RagProductManagement.d4469a80')">
          <el-input 
            v-model="form.description" 
            type="textarea"
            :rows="3"
            :placeholder="$t('auto.RagProductManagement.43d4d58b')"
          ></el-input>
        </el-form-item>
        
        <el-form-item :label="$t('auto.RagProductManagement.ffecc9f0')">
          <el-switch
            v-model="form.is_active"
            :active-text="$t('auto.RagProductManagement.83a991d7')"
            :inactive-text="$t('auto.RagProductManagement.d70e9bdf')"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('auto.RagProductManagement.625fb26b') }}</el-button>
          <el-button type="primary" @click="submitForm">{{ $t('auto.RagProductManagement.38cf16f2') }}</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ragProductConfigAPI } from '@/api/ragProductConfig'

const { t } = useI18n()

const products = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const formRef = ref()

const form = reactive({
  id: '',
  name: '',
  category: '',
  llm_model: 'gpt-3.5-turbo',
  temperature: 0.7,
  max_tokens: 1000,
  top_p: 0.9,
  frequency_penalty: 0.5,
  presence_penalty: 0.5,
  description: '',
  is_active: true,
  llm_provider_config: { base_url: '', api_key: '', model: '', api_type: 'openai' },
  embedding_provider_config: { base_url: '', api_key: '', model: '', dimension: 1024, enabled: true, api_type: 'openai' },
  rerank_provider_config: { base_url: '', api_key: '', model: '', enabled: true, api_type: 'openai' }
})

const rules = {
  name: [
    { required: true, message: i18n.global.t(t('auto.RagProductManagement.448ab9f6')), trigger: 'blur' },
    { min: 1, max: 50, message: i18n.global.t(t('auto.RagProductManagement.9d2c3eb1')), trigger: 'blur' }
  ],
  category: [
    { required: true, message: i18n.global.t(t('auto.RagProductManagement.bd663c8f')), trigger: 'change' }
  ],
  temperature: [
    { required: true, message: i18n.global.t(t('auto.RagProductManagement.c82179bc')), trigger: 'blur' },
    { type: 'number', min: 0, max: 2, message: i18n.global.t(t('auto.RagProductManagement.a7261fc0')), trigger: 'blur' }
  ],
  max_tokens: [
    { required: true, message: i18n.global.t(t('auto.RagProductManagement.1478c29c')), trigger: 'blur' },
    { type: 'number', min: 1, max: 4096, message: i18n.global.t(t('auto.RagProductManagement.8c549316')), trigger: 'blur' }
  ],
  top_p: [
    { required: true, message: i18n.global.t(t('auto.RagProductManagement.eb8cd2cd')), trigger: 'blur' },
    { type: 'number', min: 0, max: 1, message: i18n.global.t('Top-P值必须在 0 到 1 之间'), trigger: 'blur' }
  ],
  frequency_penalty: [
    { required: true, message: i18n.global.t(t('auto.RagProductManagement.85fd7ab0')), trigger: 'blur' },
    { type: 'number', min: -2, max: 2, message: i18n.global.t(t('auto.RagProductManagement.b53b46a4')), trigger: 'blur' }
  ],
  presence_penalty: [
    { required: true, message: i18n.global.t(t('auto.RagProductManagement.97a4f99a')), trigger: 'blur' },
    { type: 'number', min: -2, max: 2, message: i18n.global.t(t('auto.RagProductManagement.f08ec47c')), trigger: 'blur' }
  ]
}

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const dialogTitle = computed(() => {
  return form.id ? t('auto.RagProductManagement.0098e149') : t('auto.RagProductManagement.127cb7d7')
})

onMounted(() => {
  loadProducts()
})

const CATEGORY_LABELS = {
  agent_persona: t('auto.RagProductManagement.a7cb5642'),
  sales_script: t('auto.RagProductManagement.b7232614'),
  knowledge: t('auto.RagProductManagement.7e433388'),
  prompt: t('auto.RagProductManagement.47b7af95'),
  workflow: t('auto.RagProductManagement.228526dc'),
  tool: t('auto.RagProductManagement.20dce2c6')
};
const getCategoryLabel = (c) => CATEGORY_LABELS[c] || c || '-'
const getCategoryTagType = (category) => {
  if (!category)
    return 'info';

  const types = ['primary', 'success', 'warning', 'danger', 'info']
  let hash = 0
  for (let i = 0; i < category.length; i++) {
    const char = category.charCodeAt(i)
    hash = ((hash << 5) - hash) + char
    hash = hash & hash;
  }
  const index = Math.abs(hash) % types.length
  return types[index]
}

const loadProducts = async () => {
  loading.value = true
  try {
    const response = await ragProductConfigAPI.getRagProducts({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    products.value = response.items || [];
    pagination.total = response.total || 0
  } catch (error) {
    ElMessage.error(t('auto.RagProductManagement.24b6c3d3') + error.message)
  } finally {
    loading.value = false
  }
}

const showCreateDialog = () => {
  resetForm()
  dialogVisible.value = true
}

const editProduct = (product) => {
  Object.assign(form, product)
  form.temperature = product.temperature ?? 0.7;
  form.max_tokens = product.max_tokens ?? 1000
  form.top_p = product.top_p ?? 0.9
  form.frequency_penalty = product.frequency_penalty ?? 0.5
  form.presence_penalty = product.presence_penalty ?? 0.5
  dialogVisible.value = true
}

const toggleActive = async (product) => {
  try {
    await ElMessageBox.confirm(
      `确认${product.is_active ? t('auto.RagProductManagement.5c56a889') : t('auto.RagProductManagement.83a991d7')}该RAG产品吗？`,
      t('auto.RagProductManagement.02d9819d'),
      {
        confirmButtonText: t('auto.RagProductManagement.38cf16f2'),
        cancelButtonText: t('auto.RagProductManagement.625fb26b'),
        type: 'warning'
      }
    )
    
    const updatedProduct = {
      ...product,
      is_active: !product.is_active
    }
    
    await ragProductConfigAPI.updateRagProduct(product.id, updatedProduct)
    ElMessage.success(`${product.is_active ? t('auto.RagProductManagement.5c56a889') : t('auto.RagProductManagement.83a991d7')}成功`)
    loadProducts()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('auto.RagProductManagement.be4a600c') + error.message)
    }
  }
}

const deleteProduct = async (id) => {
  try {
    await ElMessageBox.confirm(t('auto.RagProductManagement.ef5cba3a'), t('auto.RagProductManagement.02d9819d'), {
      confirmButtonText: t('auto.RagProductManagement.38cf16f2'),
      cancelButtonText: t('auto.RagProductManagement.625fb26b'),
      type: 'warning'
    })
    
    await ragProductConfigAPI.deleteRagProduct(id)
    ElMessage.success(i18n.global.t(t('auto.RagProductManagement.0007d170')))
    loadProducts()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('auto.RagProductManagement.ad23f072') + error.message)
    }
  }
}

const resetForm = () => {
  Object.assign(form, {
    id: '',
    name: '',
    category: '',
    llm_model: 'gpt-3.5-turbo',
    temperature: 0.7,
    max_tokens: 1000,
    top_p: 0.9,
    frequency_penalty: 0.5,
    presence_penalty: 0.5,
    description: '',
    is_active: true,
    llm_provider_config: { base_url: '', api_key: '', model: '', api_type: 'openai' },
    embedding_provider_config: { base_url: '', api_key: '', model: '', dimension: 1024, enabled: true, api_type: 'openai' },
    rerank_provider_config: { base_url: '', api_key: '', model: '', enabled: true, api_type: 'openai' }
  })
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

const submitForm = async () => {
  try {
    await formRef.value.validate()
    
    if (form.id) {
      await ragProductConfigAPI.updateRagProduct(form.id, { ...form });
      ElMessage.success(i18n.global.t(t('auto.RagProductManagement.55aa6366')))
    } else {
      await ragProductConfigAPI.createRagProduct({ ...form });
      ElMessage.success(i18n.global.t(t('auto.RagProductManagement.04a691b3')))
    }
    
    dialogVisible.value = false
    loadProducts()
  } catch (error) {
    if (error?.message) {
      ElMessage.error(t('auto.RagProductManagement.931ce4df') + error.message)
    }
  }
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  loadProducts()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  loadProducts()
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>