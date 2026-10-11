<template>
  <div class="sop-template-editor-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ isEdit ? '编辑 SOP 模板' : '创建 SOP 模板' }}</h2>
          <p class="subtitle" v-if="isEdit">模板 ID：{{ templateId }} · 修改后保存生效</p>
          <p class="subtitle" v-else>{{ $t('auto.Editor.7eff0774') }}</p>
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
            <el-icon><Tickets /></el-icon>
            模板信息
          </div>
        </template>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('auto.Editor.a5d1c511')" prop="name">
              <el-input v-model="form.name" :placeholder="$t('auto.Editor.1dc18383')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Editor.764aca4e')" prop="intent">
              <el-select
                v-model="form.intent"
                filterable
                allow-create
                clearable
                :placeholder="$t('auto.Editor.3814c5cd')"
                style="width: 100%"
              >
                <el-option v-for="i in presetIntents" :key="i" :label="i" :value="i" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Editor.d316330b')" prop="stage">
              <el-select v-model="form.stage" :placeholder="$t('auto.Editor.c25adff4')" style="width: 100%">
                <el-option v-for="s in stageOptions" :key="s.value" :label="s.label" :value="s.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Editor.ee8ecb9e')">
              <el-input-number v-model="form.priority" :min="-100" :max="100" :step="1" style="width: 100%" />
              <div class="form-tip">{{ $t('auto.Editor.7a8f8b16') }}</div>
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      
      <el-card shadow="never" class="section-card">
        <template #header>
          <div class="card-title">
            <el-icon><Document /></el-icon>
            模板内容（Go text/template）
          </div>
        </template>
        <el-form-item :label="$t('auto.Editor.e5e2fd1c')" prop="template">
          <el-input
            v-model="form.template"
            type="textarea"
            :rows="6"
            placeholder="如：您好 {{.customer_name}}，关于您咨询的 {{.product_name}}，我们将在 {{.shipping_days}} 天内发货。"
          />
          <div class="form-tip">
            变量使用 <code v-pre>{{.var_name}}</code> 格式；保存后会在 Layer1 命中时用真实变量替换
          </div>
        </el-form-item>
        <el-form-item :label="$t('auto.Editor.5fec0da4')">
          <el-input
            v-model="form.vars"
            type="textarea"
            :rows="3"
            placeholder='JSON 格式：{"var_name": {"desc": "描述", "example": "示例"}}'
          />
          <div class="form-tip">{{ $t('auto.Editor.ff58df81') }}</div>
        </el-form-item>
        <el-form-item v-if="varsList.length" :label="$t('auto.Editor.ebf8e213')">
          <div class="var-preview">
            <el-tag v-for="v in varsList" :key="v.key" size="small" type="info" effect="plain" class="var-tag">
              <span class="var-name" v-pre>{{.</span><span class="var-name">{{ v.key }}</span><span class="var-name" v-pre>}}</span>
              <span class="var-desc">{{ v.desc || v.example || v.key }}</span>
            </el-tag>
          </div>
        </el-form-item>
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
            <el-form-item :label="$t('auto.Editor.8f3599b6')">
              <el-slider
                v-model="form.confidence"
                :min="0"
                :max="1"
                :step="0.05"
                show-input
                :show-input-controls="false"
              />
              <div class="form-tip">{{ $t('auto.Editor.a2a69cf8') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Editor.7854b52a')">
              <el-switch v-model="form.enabled" />
              <div class="form-tip">{{ $t('auto.Editor.a1e51532') }}</div>
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      
      <div class="footer-actions">
        <el-button @click="goBack">{{ $t('auto.Editor.625fb26b') }}</el-button>
        <el-button type="primary" :loading="saving" @click="onSave">
          <el-icon><Check /></el-icon>
          {{ isEdit ? '保存更新' : '创建模板' }}
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
import { Tickets, Document, Setting, Check } from '@element-plus/icons-vue'
import { sopTemplateApi } from '@/api/sopTemplate'

const t = i18n.global.t
const route = useRoute()
const router = useRouter()

const formRef = ref();
const pageLoading = ref(false)
const saving = ref(false)
const templateId = computed(() => route.params.id)
const isEdit = computed(() => !!templateId.value)

const presetIntents = [
  'shipping', 'refund', 'aftersales', 'pricing', 'product', 'payment', 'greeting', 'closing'
];

const stageOptions = [
  { value: 'initial', label: 'initial · 初始' },
  { value: 'middle', label: 'middle · 推进' },
  { value: 'late', label: 'late · 后期' },
  { value: 'objection', label: 'objection · 异议' },
  { value: 'closing', label: 'closing · 收单' }
];

const getDefaultForm = () => ({
  name: '',
  intent: '',
  stage: 'middle',
  template: '',
  vars: '',
  priority: 0,
  confidence: 0.8,
  enabled: true
});

const form = reactive(getDefaultForm())

const rules = {
  name: [
    { required: true, message: t(t('auto.Editor.8f21b986')), trigger: 'blur' }
  ],
  intent: [
    { required: true, message: t(t('auto.Editor.dfc9bf0b')), trigger: 'change' }
  ],
  stage: [
    { required: true, message: t(t('auto.Editor.a9fe8847')), trigger: 'change' }
  ],
  template: [
    { required: true, message: t(t('auto.Editor.73536b73')), trigger: 'blur' }
  ]
};

const varsList = computed(() => {
  if (!form.vars) return []
  try {
    const obj = JSON.parse(form.vars)
    return Object.keys(obj).map((k) => ({
      key: k,
      desc: obj[k]?.desc || '',
      example: obj[k]?.example || ''
    }))
  } catch (e) {
    return []
  }
});

const loadDetail = async () => {
  if (!isEdit.value) return
  pageLoading.value = true
  try {
    const res = await sopTemplateApi.get(templateId.value)
    if (res) {
      Object.assign(form, getDefaultForm(), res)
      form.enabled = res.enabled !== false
    }
  } catch (e) {
    ElMessage.error(t('auto.Editor.0dfd38c9') + (e.message || t('auto.Editor.974e7484')))
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
    saving.value = true
    try {
      const data = {
        name: form.name,
        intent: form.intent,
        stage: form.stage,
        template: form.template,
        vars: form.vars || '',
        priority: form.priority || 0,
        confidence: form.confidence || 0,
        enabled: form.enabled
      }
      if (isEdit.value) {
        await sopTemplateApi.update(templateId.value, data)
        ElMessage.success(t(t('auto.Editor.b846b32f')))
      } else {
        await sopTemplateApi.create(data)
        ElMessage.success(t(t('auto.Editor.49c62b01')))
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
  router.push('/sop-template/list')
}

onMounted(() => {
  loadDetail()
})
</script>

<style scoped>
.sop-template-editor-page {
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

.form-tip code {
  background: var(--el-fill-color-light);
  padding: 1px 4px;
  border-radius: 3px;
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 11px;
  color: var(--el-color-primary);
}

.var-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding: 4px 0;
}

.var-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.var-name {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-weight: 600;
  color: var(--el-color-primary);
}

.var-desc {
  color: var(--el-text-color-secondary);
  font-size: 11px;
}

.footer-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 16px 0;
}
</style>
