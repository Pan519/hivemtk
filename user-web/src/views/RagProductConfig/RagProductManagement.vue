<template>
  <div class="rag-product-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.RagProductManagement.31.99d896') }}</span>
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
        <el-table-column prop="name" :label="$t('auto.RagProductManagement.1.ff6855')" width="150"></el-table-column>
        <el-table-column prop="category" :label="$t('auto.RagProductManagement.2.383a6d')" width="100">
          <template #default="{ row }">
            <el-tag :type="getCategoryTagType(row.category)">
              {{ getCategoryLabel(row.category) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="llm_model" :label="$t('auto.RagProductManagement.3.248d0f')" width="150"></el-table-column>
        <el-table-column prop="temperature" :label="$t('auto.RagProductManagement.4.723e54')" width="100">
          <template #default="{ row }">
            {{ row.temperature.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="max_tokens" :label="$t('auto.RagProductManagement.5.b20b40')" width="100">
          <template #default="{ row }">
            {{ row.max_tokens }}
          </template>
        </el-table-column>
        <el-table-column prop="top_p" label="Top-P" width="100">
          <template #default="{ row }">
            {{ row.top_p.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="description" :label="$t('auto.RagProductManagement.6.3bdd08')" show-overflow-tooltip></el-table-column>
        <el-table-column :label="$t('auto.RagProductManagement.7.3fea7c')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_active ? 'success' : 'info'">
              {{ row.is_active ? '激活' : '未激活' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.RagProductManagement.8.eca37c')" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.RagProductManagement.9.2b6bc0')" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="editProduct(row)">{{ $t('auto.RagProductManagement.32.95b351') }}</el-button>
            <el-button 
              size="small" 
              :type="row.is_active ? 'warning' : 'success'"
              @click="toggleActive(row)"
            >
              {{ row.is_active ? '停用' : '激活' }}
            </el-button>
            <el-button size="small" type="danger" @click="deleteProduct(row.id)">{{ $t('auto.RagProductManagement.33.2f4aad') }}</el-button>
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
        <el-form-item :label="$t('auto.RagProductManagement.1.ff6855')" prop="name">
          <el-input 
            v-model="form.name" 
            :placeholder="$t('auto.RagProductManagement.10.448ab9')"
            :disabled="!!form.id"
          ></el-input>
        </el-form-item>
        
        <el-form-item :label="$t('auto.RagProductManagement.11.8d5425')" prop="category">
          <el-input 
            v-model="form.category" 
            :placeholder="$t('auto.RagProductManagement.12.0bae2a')"
            style="width: 100%;"
          ></el-input>
        </el-form-item>
        
        <el-form-item :label="$t('auto.RagProductManagement.3.248d0f')">
          <el-input 
            v-model="form.llm_model" 
            :placeholder="$t('auto.RagProductManagement.13.4b2878')"
          ></el-input>
        </el-form-item>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('auto.RagProductManagement.4.723e54')" prop="temperature">
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
            <el-form-item :label="$t('auto.RagProductManagement.14.19c9a1')" prop="top_p">
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
            <el-form-item :label="$t('auto.RagProductManagement.15.195143')" prop="max_tokens">
              <el-input-number
                v-model="form.max_tokens"
                :min="1"
                :max="4096"
                style="width: 100%;"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.RagProductManagement.16.b8f0c7')" prop="frequency_penalty">
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
        
        <el-form-item :label="$t('auto.RagProductManagement.17.ca5ae4')" prop="presence_penalty">
          <el-slider
            v-model="form.presence_penalty"
            :min="-2"
            :max="2"
            :step="0.01"
            show-input
          />
        </el-form-item>
        
        <el-divider content-position="left">{{ $t('auto.RagProductManagement.34.bd540b') }}</el-divider>

        
        <el-divider content-position="left">{{ $t('auto.RagProductManagement.35.13ada8') }}</el-divider>
        <el-form-item label="API BaseURL">
          <el-input v-model="form.llm_provider_config.base_url" :placeholder="$t('auto.RagProductManagement.18.dd75d3')" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="form.llm_provider_config.api_key" type="password" show-password placeholder="sk-..." />
        </el-form-item>
        <el-form-item :label="$t('auto.RagProductManagement.19.8000f1')">
          <el-input v-model="form.llm_provider_config.model" :placeholder="$t('auto.RagProductManagement.20.b67062')" />
        </el-form-item>

        
        <el-divider content-position="left">{{ $t('auto.RagProductManagement.36.830192') }}</el-divider>
        <el-form-item label="API BaseURL">
          <el-input v-model="form.embedding_provider_config.base_url" :placeholder="$t('auto.RagProductManagement.21.7a18dd')" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="form.embedding_provider_config.api_key" type="password" show-password placeholder="sk-..." />
        </el-form-item>
        <el-form-item :label="$t('auto.RagProductManagement.19.8000f1')">
          <el-input v-model="form.embedding_provider_config.model" :placeholder="$t('auto.RagProductManagement.22.8a31b1')" />
        </el-form-item>
        <el-form-item :label="$t('auto.RagProductManagement.23.f29c54')">
          <el-input-number v-model="form.embedding_provider_config.dimension" :min="256" :max="4096" :step="128" />
        </el-form-item>
        <el-form-item :label="$t('auto.RagProductManagement.24.7854b5')">
          <el-switch v-model="form.embedding_provider_config.enabled" />
        </el-form-item>

        
        <el-divider content-position="left">{{ $t('auto.RagProductManagement.37.980d3e') }}</el-divider>
        <el-form-item label="API BaseURL">
          <el-input v-model="form.rerank_provider_config.base_url" :placeholder="$t('auto.RagProductManagement.21.7a18dd')" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="form.rerank_provider_config.api_key" type="password" show-password placeholder="sk-..." />
        </el-form-item>
        <el-form-item :label="$t('auto.RagProductManagement.19.8000f1')">
          <el-input v-model="form.rerank_provider_config.model" :placeholder="$t('auto.RagProductManagement.25.742fab')" />
        </el-form-item>
        <el-form-item :label="$t('auto.RagProductManagement.24.7854b5')">
          <el-switch v-model="form.rerank_provider_config.enabled" />
        </el-form-item>

        <el-form-item :label="$t('auto.RagProductManagement.26.d4469a')">
          <el-input 
            v-model="form.description" 
            type="textarea"
            :rows="3"
            :placeholder="$t('auto.RagProductManagement.27.43d4d5')"
          ></el-input>
        </el-form-item>
        
        <el-form-item :label="$t('auto.RagProductManagement.28.ffecc9')">
          <el-switch
            v-model="form.is_active"
            :active-text="$t('auto.RagProductManagement.29.83a991')"
            :inactive-text="$t('auto.RagProductManagement.30.d70e9b')"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('auto.RagProductManagement.38.625fb2') }}</el-button>
          <el-button type="primary" @click="submitForm">{{ $t('auto.RagProductManagement.39.38cf16') }}</el-button>
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
    { required: true, message: i18n.global.t(t('auto.RagProductManagement.10.448ab9')), trigger: 'blur' },
    { min: 1, max: 50, message: i18n.global.t(t('auto.RagProductManagement.40.9d2c3e')), trigger: 'blur' }
  ],
  category: [
    { required: true, message: i18n.global.t(t('auto.RagProductManagement.41.bd663c')), trigger: 'change' }
  ],
  temperature: [
    { required: true, message: i18n.global.t(t('auto.RagProductManagement.42.c82179')), trigger: 'blur' },
    { type: 'number', min: 0, max: 2, message: i18n.global.t(t('auto.RagProductManagement.43.a7261f')), trigger: 'blur' }
  ],
  max_tokens: [
    { required: true, message: i18n.global.t(t('auto.RagProductManagement.44.1478c2')), trigger: 'blur' },
    { type: 'number', min: 1, max: 4096, message: i18n.global.t(t('auto.RagProductManagement.45.8c5493')), trigger: 'blur' }
  ],
  top_p: [
    { required: true, message: i18n.global.t(t('auto.RagProductManagement.46.eb8cd2')), trigger: 'blur' },
    { type: 'number', min: 0, max: 1, message: i18n.global.t('Top-P值必须在 0 到 1 之间'), trigger: 'blur' }
  ],
  frequency_penalty: [
    { required: true, message: i18n.global.t(t('auto.RagProductManagement.47.85fd7a')), trigger: 'blur' },
    { type: 'number', min: -2, max: 2, message: i18n.global.t(t('auto.RagProductManagement.48.b53b46')), trigger: 'blur' }
  ],
  presence_penalty: [
    { required: true, message: i18n.global.t(t('auto.RagProductManagement.49.97a4f9')), trigger: 'blur' },
    { type: 'number', min: -2, max: 2, message: i18n.global.t(t('auto.RagProductManagement.50.f08ec4')), trigger: 'blur' }
  ]
}

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const dialogTitle = computed(() => {
  return form.id ? t('auto.RagProductManagement.51.0098e1') : t('auto.RagProductManagement.52.127cb7')
})

onMounted(() => {
  loadProducts()
})

const CATEGORY_LABELS = {
  agent_persona: t('auto.RagProductManagement.53.a7cb56'),
  sales_script: t('auto.RagProductManagement.54.b72326'),
  knowledge: t('auto.RagProductManagement.55.7e4333'),
  prompt: t('auto.RagProductManagement.56.47b7af'),
  workflow: t('auto.RagProductManagement.57.228526'),
  tool: t('auto.RagProductManagement.58.20dce2')
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
    ElMessage.error(t('auto.RagProductManagement.59.24b6c3') + error.message)
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
      `确认${product.is_active ? t('auto.RagProductManagement.60.5c56a8') : t('auto.RagProductManagement.29.83a991')}该RAG产品吗？`,
      t('auto.RagProductManagement.61.02d981'),
      {
        confirmButtonText: t('auto.RagProductManagement.39.38cf16'),
        cancelButtonText: t('auto.RagProductManagement.38.625fb2'),
        type: 'warning'
      }
    )
    
    const updatedProduct = {
      ...product,
      is_active: !product.is_active
    }
    
    await ragProductConfigAPI.updateRagProduct(product.id, updatedProduct)
    ElMessage.success(`${product.is_active ? t('auto.RagProductManagement.60.5c56a8') : t('auto.RagProductManagement.29.83a991')}成功`)
    loadProducts()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('auto.RagProductManagement.62.be4a60') + error.message)
    }
  }
}

const deleteProduct = async (id) => {
  try {
    await ElMessageBox.confirm(t('auto.RagProductManagement.63.ef5cba'), t('auto.RagProductManagement.61.02d981'), {
      confirmButtonText: t('auto.RagProductManagement.39.38cf16'),
      cancelButtonText: t('auto.RagProductManagement.38.625fb2'),
      type: 'warning'
    })
    
    await ragProductConfigAPI.deleteRagProduct(id)
    ElMessage.success(i18n.global.t(t('auto.RagProductManagement.64.0007d1')))
    loadProducts()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('auto.RagProductManagement.65.ad23f0') + error.message)
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
      ElMessage.success(i18n.global.t(t('auto.RagProductManagement.66.55aa63')))
    } else {
      await ragProductConfigAPI.createRagProduct({ ...form });
      ElMessage.success(i18n.global.t(t('auto.RagProductManagement.67.04a691')))
    }
    
    dialogVisible.value = false
    loadProducts()
  } catch (error) {
    if (error?.message) {
      ElMessage.error(t('auto.RagProductManagement.68.931ce4') + error.message)
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