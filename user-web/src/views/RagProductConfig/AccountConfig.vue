<template>
  <div class="account-config">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('账号配置') }}</span>
        </div>
      </template>
      
      <el-form 
        :model="form" 
        :rules="rules" 
        ref="formRef"
        label-width="120px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('auto.AccountConfig.4981287f')" prop="platform">
              <el-select
                v-model="form.platform"
                :placeholder="$t('auto.AccountConfig.9134b450')"
                @change="onPlatformChange"
                style="width: 100%;"
              >
                <el-option
                  v-for="opt in cardPlatformOptions"
                  :key="opt.value"
                  :label="opt.label"
                  :value="opt.value"
                />
              </el-select>
            </el-form-item>
          </el-col>
          
          <el-col :span="12">
            <el-form-item :label="$t('auto.AccountConfig.14d1718c')" prop="account_id">
              <el-input 
                v-model="form.account_id" 
                :placeholder="$t('auto.AccountConfig.c03276fd')"
              ></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('auto.AccountConfig.6e7ce48d')">
              <el-select 
                v-model="form.rag_product_id" 
                :placeholder="$t('auto.AccountConfig.85e9a88b')"
                :disabled="!ragProducts.length"
                style="width: 100%;"
              >
                <el-option 
                  v-for="product in ragProducts" 
                  :key="product.id" 
                  :label="product.name" 
                  :value="product.id"
                >
                  <span>{{ product.name }}</span>
                  <el-tag 
                    size="small" 
                    type="info" 
                    style="margin-left: 8px;"
                  >
                    {{ product.category }}
                  </el-tag>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          
          <el-col :span="12">
            <el-form-item :label="$t('auto.AccountConfig.ffb248fd')">
              <el-input-number 
                v-model="form.max_daily_queries" 
                :min="1" 
                :max="10000"
                style="width: 100%;"
              ></el-input-number>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item :label="$t('auto.AccountConfig.b9613078')">
          <el-switch
            v-model="form.is_auto_reply_enabled"
            :active-text="$t('auto.AccountConfig.922a7731')"
            style="margin-right: 20px;"
          />
          <el-switch
            v-model="form.is_rag_enabled"
            :active-text="$t('auto.AccountConfig.cb093b41')"
            :disabled="!form.rag_product_id"
          />
        </el-form-item>
        
        <el-form-item :label="$t('auto.AccountConfig.988f37d4')">
          <el-button 
            type="primary" 
            size="small" 
            @click="addReplyRule"
          >
            添加规则
          </el-button>
        </el-form-item>
        
        <el-form-item>
          <div 
            v-for="(rule, index) in form.reply_rules" 
            :key="index"
            class="reply-rule-item"
          >
            <el-card shadow="hover">
              <template #header>
                <div class="rule-header">
                  <span>规则 {{ index + 1 }}</span>
                  <el-button 
                    type="danger" 
                    size="small" 
                    @click="removeReplyRule(index)"
                  >
                    删除
                  </el-button>
                </div>
              </template>
              
              <el-row :gutter="10">
                <el-col :span="6">
                  <el-input 
                    v-model.number="rule.priority" 
                    :placeholder="$t('auto.AccountConfig.ee8ecb9e')"
                    type="number"
                  ></el-input>
                </el-col>
                
                <el-col :span="12">
                  <el-input 
                    v-model="rule.keywords_str" 
                    :placeholder="$t('auto.AccountConfig.f42e6971')"
                  ></el-input>
                </el-col>
                
                <el-col :span="6">
                  <el-switch
                    v-model="rule.is_active"
                    :active-text="$t('auto.AccountConfig.7854b52a')"
                  />
                </el-col>
              </el-row>
              
              <el-input 
                v-model="rule.reply_template" 
                :placeholder="$t('auto.AccountConfig.0efefdeb')"
                type="textarea"
                :rows="2"
                style="margin-top: 10px;"
              ></el-input>
            </el-card>
          </div>
        </el-form-item>
        
        <el-form-item>
          <el-button 
            type="primary" 
            @click="submitForm"
          >
            保存配置
          </el-button>
          <el-button @click="resetForm">{{ $t('auto.AccountConfig.4b9c3271') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { ragProductConfigAPI } from '@/api/ragProductConfig'
import { CARD_PLATFORM_OPTIONS } from '@/constants/cardPlatform';

const cardPlatformOptions = CARD_PLATFORM_OPTIONS

const { t } = useI18n()

const formRef = ref()

const form = reactive({
  platform: '',
  account_id: '',
  rag_product_id: '',
  is_auto_reply_enabled: false,
  is_rag_enabled: false,
  max_daily_queries: 1000,
  reply_rules: []
})

const rules = {
  platform: [
    { required: true, message: i18n.global.t(t('auto.AccountConfig.72b562ed')), trigger: 'change' }
  ],
  account_id: [
    { required: true, message: i18n.global.t(t('auto.AccountConfig.cf5e16f7')), trigger: 'blur' }
  ]
}

const ragProducts = ref([])
const currentConfig = ref(null)

onMounted(() => {
  loadRagProducts()
})

const loadRagProducts = async () => {
  try {
    const response = await ragProductConfigAPI.getRagProducts()
    ragProducts.value = response.items || [];
  } catch (error) {
    ElMessage.error(t('auto.AccountConfig.24b6c3d3') + error.message)
  }
}

const onPlatformChange = () => {
  form.account_id = '';
  loadCurrentConfig()
}

const loadCurrentConfig = async () => {
  if (!form.account_id || !form.platform) return
  
  try {
    const response = await ragProductConfigAPI.getAccountConfig({
      account_id: form.account_id,
      platform: form.platform
    })
    
    const config = response;
    Object.assign(form, {
      ...config,
      rag_product_id: config.rag_product_id || '',
      reply_rules: (config.reply_rules || []).map(rule => ({
        ...rule,
        keywords_str: Array.isArray(rule.keywords) ? rule.keywords.join(',') : rule.keywords || ''
      }))
    })
    currentConfig.value = config
  } catch (error) {
    form.reply_rules = [];
    currentConfig.value = null
  }
}

const addReplyRule = () => {
  form.reply_rules.push({
    id: Date.now().toString(),
    priority: 1,
    keywords_str: '',
    reply_template: '',
    is_active: true
  })
}

const removeReplyRule = (index) => {
  form.reply_rules.splice(index, 1)
}

const submitForm = async () => {
  try {
    await formRef.value.validate()
    
    const submitData = {
      ...form,
      reply_rules: form.reply_rules.map(rule => ({
        ...rule,
        keywords: rule.keywords_str ? rule.keywords_str.split(',').map(k => k.trim()).filter(k => k) : []
      }))
    }
    
    await ragProductConfigAPI.updateAccountConfig(submitData)
    ElMessage.success(i18n.global.t(t('auto.AccountConfig.0a4915ab')))
    
    loadCurrentConfig();
  } catch (error) {
    ElMessage.error(t('auto.AccountConfig.40f90217') + error.message)
  }
}

const resetForm = () => {
  formRef.value.resetFields()
  if (currentConfig.value) {
    Object.assign(form, currentConfig.value)
  } else {
    form.reply_rules = []
  }
}
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.reply-rule-item {
  margin-bottom: 15px;
}

.rule-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
