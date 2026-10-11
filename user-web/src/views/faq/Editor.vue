<template>
  <div class="faq-editor-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ isEdit ? '编辑 FAQ' : '创建 FAQ' }}</h2>
          <p class="subtitle" v-if="isEdit">FAQ ID：{{ faqId }} · 修改后保存生效</p>
          <p class="subtitle" v-else>{{ $t('auto.Editor.14.8ae443') }}</p>
        </div>
        <div class="header-actions">
          <el-button @click="goBack">{{ $t('auto.Editor.15.adcd1d') }}</el-button>
        </div>
      </div>
    </el-card>

    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
      v-loading="pageLoading"
    >
      
      <el-card shadow="never" class="section-card">
        <template #header>
          <div class="card-title">
            <el-icon><QuestionFilled /></el-icon>
            FAQ 内容
          </div>
        </template>
        <el-form-item :label="$t('auto.Editor.1.5dc99f')" prop="question">
          <el-input
            v-model="form.question"
            type="textarea"
            :rows="2"
            :placeholder="$t('auto.Editor.2.40bae0')"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.Editor.3.dbacf3')" prop="answer">
          <el-input
            v-model="form.answer"
            type="textarea"
            :rows="4"
            :placeholder="$t('auto.Editor.4.fddf8d')"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.Editor.5.9699a5')" prop="keywords">
          <el-input
            v-model="form.keywordsInput"
            :placeholder="$t('auto.Editor.6.9667af')"
          />
          <div class="form-tip">{{ $t('auto.Editor.16.20986a') }}</div>
        </el-form-item>
        <el-form-item v-if="form.keywords && form.keywords.length" :label="$t('auto.Editor.7.645dbc')">
          <div class="keyword-preview">
            <el-tag v-for="k in form.keywords" :key="k" size="small" type="info" effect="plain" closable @close="removeKeyword(k)">
              {{ k }}
            </el-tag>
          </div>
        </el-form-item>
      </el-card>

      
      <el-card shadow="never" class="section-card">
        <template #header>
          <div class="card-title">
            <el-icon><CollectionTag /></el-icon>
            分类与意图
          </div>
        </template>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('auto.Editor.8.e50e8a')">
              <el-select
                v-model="form.category"
                filterable
                allow-create
                clearable
                :placeholder="$t('auto.Editor.9.704203')"
                style="width: 100%"
              >
                <el-option v-for="c in presetCategories" :key="c" :label="c" :value="c" />
              </el-select>
              <div class="form-tip">{{ $t('auto.Editor.17.4de8f8') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Editor.10.852551')">
              <el-input v-model="form.intent" :placeholder="$t('auto.Editor.11.f56009')" />
              <div class="form-tip">{{ $t('auto.Editor.18.231c22') }}</div>
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      
      <el-card shadow="never" class="section-card">
        <template #header>
          <div class="card-title">
            <el-icon><Setting /></el-icon>
            匹配参数
          </div>
        </template>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('auto.Editor.12.44c611')">
              <el-slider
                v-model="form.confidence"
                :min="0"
                :max="1"
                :step="0.05"
                show-input
                :show-input-controls="false"
              />
              <div class="form-tip">{{ $t('auto.Editor.19.a35e18') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Editor.13.7854b5')">
              <el-switch v-model="form.enabled" />
              <div class="form-tip">{{ $t('auto.Editor.20.bfd59d') }}</div>
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      
      <div class="footer-actions">
        <el-button @click="goBack">{{ $t('auto.Editor.21.625fb2') }}</el-button>
        <el-button type="primary" :loading="saving" @click="onSave">
          <el-icon><Check /></el-icon>
          {{ isEdit ? '保存更新' : '创建 FAQ' }}
        </el-button>
      </div>
    </el-form>
  </div>
</template>

<script setup>
import i18n from '@/i18n'
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  QuestionFilled, CollectionTag, Setting, Check
} from '@element-plus/icons-vue'
import { faqApi } from '@/api/faq'

const t = i18n.global.t
const route = useRoute()
const router = useRouter()

const formRef = ref();
const pageLoading = ref(false)
const saving = ref(false)
const faqId = computed(() => route.params.id)
const isEdit = computed(() => !!faqId.value)

const presetCategories = [
  'logistics',
  'pricing',
  'aftersales',
  'product',
  'payment',
  'account',
  'general'
];

const getDefaultForm = () => ({
  question: '',
  answer: '',
  keywords: [],
  keywordsInput: '',
  category: '',
  intent: '',
  confidence: 0.8,
  enabled: true
});

const form = reactive(getDefaultForm())

const rules = {
  question: [
    { required: true, message: t(t('auto.Editor.22.73fb82')), trigger: 'blur' }
  ],
  answer: [
    { required: true, message: t(t('auto.Editor.23.6a34fc')), trigger: 'blur' }
  ]
};

const removeKeyword = (k) => {
  form.keywords = form.keywords.filter((x) => x !== k)
};

const parseKeywords = (input) => {
  if (!input) return []
  return input
    .split(/[,，;；\s]+/)
    .map((s) => s.trim())
    .filter(Boolean)
}

const loadDetail = async () => {
  if (!isEdit.value) return
  pageLoading.value = true
  try {
    const res = await faqApi.get(faqId.value)
    if (res) {
      Object.assign(form, getDefaultForm(), res)
      form.keywords = Array.isArray(res.keywords) ? [...res.keywords] : []
      form.keywordsInput = form.keywords.join(', ')
      form.enabled = res.enabled !== false
    }
  } catch (e) {
    ElMessage.error(t('auto.Editor.24.de5e3b') + (e.message || t('auto.Editor.25.974e74')))
    goBack()
  } finally {
    pageLoading.value = false
  }
};

const onSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) {
      ElMessage.warning(t(t('auto.Editor.26.341de0')))
      return
    }
    const fromInput = parseKeywords(form.keywordsInput);
    const merged = Array.from(new Set([...(form.keywords || []), ...fromInput]))
    saving.value = true
    try {
      const data = {
        question: form.question,
        answer: form.answer,
        keywords: merged,
        category: form.category || '',
        intent: form.intent || '',
        confidence: form.confidence || 0,
        enabled: form.enabled
      }
      if (isEdit.value) {
        await faqApi.update(faqId.value, data)
        ElMessage.success(t('FAQ 更新成功'))
      } else {
        await faqApi.create(data)
        ElMessage.success(t('FAQ 创建成功'))
      }
      goBack()
    } catch (e) {
      ElMessage.error(t('auto.Editor.27.78d235') + (e.message || t('auto.Editor.25.974e74')))
    } finally {
      saving.value = false
    }
  })
};

const goBack = () => {
  router.push('/faq/list')
}

onMounted(() => {
  loadDetail()
})
</script>

<style scoped>
.faq-editor-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.header-card :deep(.el-card__body) {
  padding: 18px 24px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}

.header-content h2 {
  margin: 0 0 4px 0;
  font-size: 20px;
  font-weight: 600;
}

.subtitle {
  margin: 0;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.section-card {
  margin-top: 0;
}

.section-card :deep(.el-card__body) {
  padding: 18px 20px 4px 20px;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.form-tip {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 4px;
  line-height: 1.4;
}

.keyword-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding: 4px 0;
}

.footer-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 16px 0;
}
</style>
