<template>
  <div class="faq-editor-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ isEdit ? '编辑 FAQ' : '创建 FAQ' }}</h2>
          <p class="subtitle" v-if="isEdit">FAQ ID：{{ faqId }} · 修改后保存生效</p>
          <p class="subtitle" v-else>{{ $t('auto.Editor.8ae4430a') }}</p>
        </div>
        <div class="header-actions">
          <el-button @click="goBack">{{ $t('auto.Editor.adcd1dd7') }}</el-button>
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
        <el-form-item :label="$t('auto.Editor.5dc99f6e')" prop="question">
          <el-input
            v-model="form.question"
            type="textarea"
            :rows="2"
            :placeholder="$t('auto.Editor.40bae09f')"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.Editor.dbacf318')" prop="answer">
          <el-input
            v-model="form.answer"
            type="textarea"
            :rows="4"
            :placeholder="$t('auto.Editor.fddf8dea')"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.Editor.9699a50e')" prop="keywords">
          <el-input
            v-model="form.keywordsInput"
            :placeholder="$t('auto.Editor.9667af9d')"
          />
          <div class="form-tip">{{ $t('auto.Editor.20986ab8') }}</div>
        </el-form-item>
        <el-form-item v-if="form.keywords && form.keywords.length" :label="$t('auto.Editor.645dbc55')">
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
            <el-form-item :label="$t('auto.Editor.e50e8a3f')">
              <el-select
                v-model="form.category"
                filterable
                allow-create
                clearable
                :placeholder="$t('auto.Editor.704203ba')"
                style="width: 100%"
              >
                <el-option v-for="c in presetCategories" :key="c" :label="c" :value="c" />
              </el-select>
              <div class="form-tip">{{ $t('auto.Editor.4de8f80e') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Editor.85255185')">
              <el-input v-model="form.intent" :placeholder="$t('auto.Editor.f5600972')" />
              <div class="form-tip">{{ $t('auto.Editor.231c2296') }}</div>
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
            <el-form-item :label="$t('auto.Editor.44c6114c')">
              <el-slider
                v-model="form.confidence"
                :min="0"
                :max="1"
                :step="0.05"
                show-input
                :show-input-controls="false"
              />
              <div class="form-tip">{{ $t('auto.Editor.a35e18ec') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Editor.7854b52a')">
              <el-switch v-model="form.enabled" />
              <div class="form-tip">{{ $t('auto.Editor.bfd59db5') }}</div>
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      
      <div class="footer-actions">
        <el-button @click="goBack">{{ $t('auto.Editor.625fb26b') }}</el-button>
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
    { required: true, message: t(t('auto.Editor.73fb8249')), trigger: 'blur' }
  ],
  answer: [
    { required: true, message: t(t('auto.Editor.6a34fc28')), trigger: 'blur' }
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
    ElMessage.error(t('auto.Editor.de5e3b6d') + (e.message || t('auto.Editor.974e7484')))
    goBack()
  } finally {
    pageLoading.value = false
  }
};

const onSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) {
      ElMessage.warning(t(t('auto.Editor.341de0e4')))
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
      ElMessage.error(t('auto.Editor.78d235cc') + (e.message || t('auto.Editor.974e7484')))
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
