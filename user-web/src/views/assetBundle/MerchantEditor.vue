<template>
  <div class="merchant-editor">
    
    <el-card class="status-bar" shadow="never">
      <div class="status-row">
        <div class="status-left">
          <el-icon v-if="form.status === 'active'" class="status-icon active"><Monitor /></el-icon>
          <el-icon v-else class="status-icon paused"><VideoPause /></el-icon>
          <div class="status-text">
            <div class="status-title">
              <span v-if="form.status === 'active'">{{ $t('auto.MerchantEditor.46.3dc4bd') }}</span>
              <span v-else>{{ $t('auto.MerchantEditor.47.e03a18') }}</span>
            </div>
            <div class="status-sub">
              资产 ID: {{ form.asset_id || '未保存' }} · 版本: {{ form.version || '1.0.0' }}
            </div>
          </div>
        </div>
        <div class="status-right">
          <el-button
            v-if="form.status === 'active'"
            type="warning"
            @click="handleToggleStatus('pause')"
          >
            ⏸️ 暂停托管
          </el-button>
          <el-button
            v-else
            type="success"
            @click="handleToggleStatus('resume')"
          >
            ▶️ 恢复托管
          </el-button>
          <el-button type="primary" :loading="saving" @click="handleSave">
            💾 保存配置并立刻热更新到商户本地 AI 引擎
          </el-button>
        </div>
      </div>
    </el-card>

    <el-form ref="formRef" :model="form" label-width="160px" class="form-body" v-loading="loading">
      
      <el-card class="section" shadow="never">
        <template #header>
          <div class="section-header">
            <el-icon><Goods /></el-icon>
            <span>{{ $t('auto.MerchantEditor.48.0f8a0c') }}</span>
            <span class="section-hint">{{ $t('auto.MerchantEditor.49.ee879b') }}</span>
          </div>
        </template>
        <el-row :gutter="20">
          <el-col :span="6">
            <el-form-item :label="$t('auto.MerchantEditor.1.641cfa')">
              <el-upload
                class="cover-uploader"
                :http-request="handleCoverUpload"
                :show-file-list="false"
                accept=".jpg,.jpeg,.png,.gif,.webp"
              >
                <el-image
                  v-if="form.cover_image"
                  :src="form.cover_image"
                  :preview-src-list="[form.cover_image]"
                  fit="cover"
                  class="cover-preview"
                />
                <el-icon v-else class="cover-uploader-icon"><Plus /></el-icon>
              </el-upload>
              <div class="cover-tip">{{ $t('auto.MerchantEditor.50.70baac') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="18">
            <el-form-item :label="$t('auto.MerchantEditor.2.0acf6d')">
              <el-input v-model="form.title" :placeholder="$t('auto.MerchantEditor.3.5a5db2')" />
            </el-form-item>
            <el-form-item :label="$t('auto.MerchantEditor.4.69847d')">
              <el-input v-model="form.asset_id" :placeholder="$t('auto.MerchantEditor.5.01e0f1')" :disabled="isEdit" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('auto.MerchantEditor.6.0ad1ab')">
              <el-input v-model="form.campaign_name" :placeholder="$t('auto.MerchantEditor.7.a7a9c1')" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item :label="$t('auto.MerchantEditor.8.59e2db')">
              <el-input v-model="form.discount_pct" :placeholder="$t('auto.MerchantEditor.9.8d8112')" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item :label="$t('auto.MerchantEditor.10.562337')">
              <el-input v-model="form.support_contact" :placeholder="$t('auto.MerchantEditor.11.241003')" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('auto.MerchantEditor.12.537325')">
              <el-input v-model="form.shop_name" :placeholder="$t('auto.MerchantEditor.13.3c02ae')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.MerchantEditor.14.be3d89')">
              <el-input v-model="form.author" :placeholder="$t('auto.MerchantEditor.15.015d6e')" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      
      <el-card class="section" shadow="never">
        <template #header>
          <div class="section-header">
            <el-icon><Lock /></el-icon>
            <span>{{ $t('auto.MerchantEditor.51.1888fa') }}</span>
            <span class="section-hint">{{ $t('auto.MerchantEditor.52.deac71') }}</span>
          </div>
        </template>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item :label="$t('auto.MerchantEditor.16.c0fbfd')">
              <el-select v-model="form.crisis_threshold" :placeholder="$t('auto.MerchantEditor.17.0a11d3')">
                <el-option :label="$t('auto.MerchantEditor.18.1b57b4')" value="2" />
                <el-option :label="$t('auto.MerchantEditor.19.1ff4bb')" value="3" />
                <el-option :label="$t('auto.MerchantEditor.20.9e6086')" value="4" />
                <el-option :label="$t('auto.MerchantEditor.21.47db26')" value="5" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.MerchantEditor.22.c947d4')">
              <el-select v-model="form.tone_level" :placeholder="$t('auto.MerchantEditor.23.8dfbf0')">
                <el-option :label="$t('auto.MerchantEditor.24.9a8e99')" value="low" />
                <el-option :label="$t('auto.MerchantEditor.25.b51f75')" value="medium" />
                <el-option :label="$t('auto.MerchantEditor.26.fd1d81')" value="high" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.MerchantEditor.27.9b4521')">
              <el-radio-group v-model="form.censorship_level">
                <el-radio label="strict">{{ $t('auto.MerchantEditor.53.15087b') }}</el-radio>
                <el-radio label="unlock">{{ $t('auto.MerchantEditor.54.50d7e2') }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('auto.MerchantEditor.28.e819d9')">
          <el-checkbox-group v-model="form.enabled_intents">
            <el-checkbox label="faq">{{ $t('auto.MerchantEditor.55.4bc10d') }}</el-checkbox>
            <el-checkbox label="lead_capture">{{ $t('auto.MerchantEditor.56.535d46') }}</el-checkbox>
            <el-checkbox label="human_transfer">{{ $t('auto.MerchantEditor.57.3e61a7') }}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-card>

      
      <el-card class="section" shadow="never">
        <template #header>
          <div class="section-header">
            <el-icon><ChatDotRound /></el-icon>
            <span>{{ $t('auto.MerchantEditor.58.08da94') }}</span>
            <span class="section-hint">{{ $t('auto.MerchantEditor.59.139adb') }}</span>
          </div>
        </template>
        <div class="qa-list">
          <el-card
            v-for="(card, idx) in form.qa_cards"
            :key="card.id"
            class="qa-card"
            shadow="hover"
          >
            <div class="qa-card-header">
              <span class="qa-card-title">▼ 行业经典答疑卡片 {{ idx + 1 }}</span>
              <el-button type="danger" link @click="removeQACard(card.id)">{{ $t('auto.MerchantEditor.60.2f4aad') }}</el-button>
            </div>
            <el-form-item :label="$t('auto.MerchantEditor.29.ae37bc')">
              <el-input v-model="card.trigger" :placeholder="$t('auto.MerchantEditor.30.02798f')" />
            </el-form-item>
            <el-form-item :label="$t('auto.MerchantEditor.31.b4baaa')">
              <el-input
                v-model="card.user_example"
                type="textarea"
                :rows="2"
                :placeholder="$t('auto.MerchantEditor.32.0b770e')"
              />
            </el-form-item>
            <el-form-item :label="$t('auto.MerchantEditor.33.9ac55b')">
              <el-input
                v-model="card.reply"
                type="textarea"
                :rows="4"
                :placeholder="$t('auto.MerchantEditor.34.9c32be')"
              />
            </el-form-item>
          </el-card>
        </div>
        <el-button type="primary" plain class="add-card-btn" @click="addQACard">
          ➕ 新增商户自家店铺专属的提问应答卡片
        </el-button>
      </el-card>

      
      <el-card class="section" shadow="never">
        <template #header>
          <div class="section-header">
            <el-icon><Grid /></el-icon>
            <span>{{ $t('auto.MerchantEditor.61.4774fd') }}</span>
            <span class="section-hint">{{ $t('auto.MerchantEditor.62.5d3663') }}</span>
          </div>
        </template>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item :label="$t('auto.MerchantEditor.35.f60f4a')">
              <el-radio-group v-model="form.card_config.intent_type">
                <el-radio label="button_card">{{ $t('auto.MerchantEditor.63.aaa40b') }}</el-radio>
                <el-radio label="coupon">{{ $t('auto.MerchantEditor.64.c323c4') }}</el-radio>
                <el-radio label="handoff">{{ $t('auto.MerchantEditor.65.658619') }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="16">
            <el-form-item :label="$t('auto.MerchantEditor.36.83747b')">
              <el-input v-model="form.card_config.product_image" :placeholder="$t('auto.MerchantEditor.37.2edb14')" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-divider>{{ $t('auto.MerchantEditor.66.4a06d4') }}</el-divider>
        <div class="button-list">
          <el-card
            v-for="(btn, idx) in form.card_config.buttons"
            :key="idx"
            class="button-card"
            shadow="hover"
          >
            <el-row :gutter="10">
              <el-col :span="8">
                <el-form-item :label="$t('auto.MerchantEditor.38.ddefac')">
                  <el-input v-model="btn.title" :placeholder="$t('auto.MerchantEditor.39.934285')" />
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item :label="$t('auto.MerchantEditor.40.db24c6')">
                  <el-select v-model="btn.action">
                    <el-option :label="$t('auto.MerchantEditor.41.e44a63')" value="open_url" />
                    <el-option :label="$t('auto.MerchantEditor.42.9e7021')" value="call_api" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item v-if="btn.action === 'open_url'" :label="$t('auto.MerchantEditor.41.e44a63')">
                  <el-input v-model="btn.url" placeholder="https://shopify.com/..." />
                </el-form-item>
                <el-form-item v-else :label="$t('auto.MerchantEditor.43.4b3dc0')">
                  <el-input v-model="btn.api_name" :placeholder="$t('auto.MerchantEditor.44.4a59e8')" />
                </el-form-item>
              </el-col>
              <el-col :span="2">
                <el-button type="danger" link @click="removeButton(idx)">{{ $t('auto.MerchantEditor.60.2f4aad') }}</el-button>
              </el-col>
            </el-row>
            <el-form-item v-if="btn.action === 'call_api'" :label="$t('auto.MerchantEditor.45.8bc2cc')">
              <el-input v-model="btn.api_args" placeholder='例如：{"coupon_id":"SUMMER15"}' />
            </el-form-item>
          </el-card>
        </div>
        <el-button type="primary" plain class="add-card-btn" @click="addButton">
          ➕ 新增动作按钮
        </el-button>
      </el-card>

      
      <div class="footer-actions">
        <el-button type="primary" size="large" :loading="saving" @click="handleSave">
          💾 保存配置并立刻热更新到商户本地 AI 引擎
        </el-button>
        <el-button size="large" @click="$router.push('/asset-bundle/list')">{{ $t('auto.MerchantEditor.67.adcd1d') }}</el-button>
      </div>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Monitor, VideoPause, Goods, Lock, ChatDotRound, Grid } from '@element-plus/icons-vue'
import { merchantSave, merchantParse, enableBundle, disableBundle, getBundleByAssetID, uploadCover } from '@/api/assetBundle'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const formRef = ref(null)
const loading = ref(false)
const saving = ref(false)

const aid = computed(() => route.params.aid || '')
const isEdit = computed(() => !!aid.value)

const form = reactive({
  asset_id: '',
  title: '',
  author: '',
  cover_image: '',
  shop_name: '',
  campaign_name: '',
  discount_pct: '',
  support_contact: '',
  crisis_threshold: '4',
  tone_level: 'medium',
  censorship_level: 'strict',
  enabled_intents: ['faq', 'lead_capture', 'human_transfer'],
  qa_cards: [],
  card_config: {
    intent_type: 'button_card',
    product_image: '',
    buttons: []
  },
  status: 'draft',
  version: '1.0.0',
  template_asset_id: ''
});

const genCardId = () => `card_${Date.now()}_${Math.floor(Math.random() * 1000)}`;

const addQACard = () => {
  form.qa_cards.push({
    id: genCardId(),
    trigger: '',
    user_example: '',
    reply: '',
    order: form.qa_cards.length
  })
}

const removeQACard = (id) => {
  const idx = form.qa_cards.findIndex(c => c.id === id)
  if (idx >= 0) form.qa_cards.splice(idx, 1)
}

const addButton = () => {
  form.card_config.buttons.push({
    title: '',
    action: 'open_url',
    url: '',
    api_name: '',
    api_args: '',
    order: form.card_config.buttons.length
  })
}

const removeButton = (idx) => {
  form.card_config.buttons.splice(idx, 1)
}

const loadBundle = async () => {
  if (!aid.value) return
  loading.value = true
  try {
    const detailResp = await getBundleByAssetID(aid.value);
    const detail = detailResp?.data || detailResp
    if (detail) {
      form.title = detail.title || ''
      form.author = detail.author || ''
      form.version = detail.version || '1.0.0'
      form.status = detail.status || 'draft'
    }
    const parseResp = await merchantParse(aid.value);
    const parsed = parseResp?.data || parseResp
    if (parsed) {
      form.shop_name = parsed.shop_name || ''
      form.campaign_name = parsed.campaign_name || ''
      form.discount_pct = parsed.discount_pct || ''
      form.support_contact = parsed.support_contact || ''
      form.qa_cards = (parsed.qa_cards || []).map(c => ({ ...c, id: c.id || genCardId() }))
      if (parsed.card_config) {
        form.card_config = {
          intent_type: parsed.card_config.intent_type || 'button_card',
          product_image: parsed.card_config.product_image || '',
          buttons: parsed.card_config.buttons || []
        }
      }
      if (parsed.enabled_intents && parsed.enabled_intents.length) {
        form.enabled_intents = parsed.enabled_intents
      }
      form.crisis_threshold = parsed.crisis_threshold || '4';
      form.tone_level = parsed.tone_level || 'medium'
      form.censorship_level = parsed.censorship_level || 'default'
    }
  } catch (e) {
    ElMessage.error(t('auto.MerchantEditor.68.2cfc14') + (e?.message || e))
  } finally {
    loading.value = false
  }
};

const handleCoverUpload = async ({ file }) => {
  try {
    const res = await uploadCover(file)
    form.cover_image = res.data?.url || res.url || res
    ElMessage.success(t('auto.MerchantEditor.69.e04a52'))
  } catch (e) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
};

const handleSave = async () => {
  if (!form.asset_id) {
    ElMessage.warning(t('auto.MerchantEditor.70.b9b9a9'))
    return
  }
  if (!form.title) {
    ElMessage.warning(t('auto.MerchantEditor.71.ac8238'))
    return
  }
  saving.value = true
  try {
    const payload = {
      asset_id: form.asset_id,
      title: form.title,
      author: form.author || userStore.username,
      shop_name: form.shop_name,
      campaign_name: form.campaign_name,
      discount_pct: form.discount_pct,
      support_contact: form.support_contact,
      crisis_threshold: form.crisis_threshold,
      tone_level: form.tone_level,
      censorship_level: form.censorship_level,
      enabled_intents: form.enabled_intents,
      qa_cards: form.qa_cards.map((c, i) => ({ ...c, order: i })),
      card_config: {
        ...form.card_config,
        buttons: form.card_config.buttons.map((b, i) => ({ ...b, order: i }))
      },
      template_asset_id: form.template_asset_id
    }
    const resp = await merchantSave(payload)
    const data = resp?.data || resp
    if (data && data.id) {
      try {
        await enableBundle(data.id)
        ElMessage.success(t('auto.MerchantEditor.72.a7f38d'))
      } catch (e) {
        ElMessage.warning(t('auto.MerchantEditor.73.045c62'))
      }
      if (!isEdit.value && data.asset_id) {
        router.replace(`/asset-bundle/merchant/${data.asset_id}`)
      }
    } else {
      ElMessage.success(t('auto.MerchantEditor.74.3b1083'))
    }
  } catch (e) {
    ElMessage.error(t('auto.MerchantEditor.75.40f902') + (e?.message || e))
  } finally {
    saving.value = false
  }
};

const handleToggleStatus = async (action) => {
  if (!aid.value) {
    ElMessage.warning(t('auto.MerchantEditor.76.9558c4'))
    return
  }
  try {
    const detailResp = await getBundleByAssetID(aid.value);
    const detail = detailResp?.data || detailResp
    if (!detail || !detail.id) {
      ElMessage.error(t('auto.MerchantEditor.77.23fe0d'))
      return
    }
    if (action === 'pause') {
      await disableBundle(detail.id)
      form.status = 'inactive'
      ElMessage.success(t('auto.MerchantEditor.78.76be2f'))
    } else {
      await enableBundle(detail.id)
      form.status = 'active'
      ElMessage.success(t('auto.MerchantEditor.79.74aa81'))
    }
  } catch (e) {
    ElMessage.error(t('auto.MerchantEditor.80.1a4a30') + (e?.message || e))
  }
};

onMounted(() => {
  if (aid.value) {
    form.asset_id = aid.value
    loadBundle()
  } else {
    addQACard();
  }
})
</script>

<style scoped>
.merchant-editor {
  padding: 16px;
}
.status-bar {
  margin-bottom: 16px;
}
.status-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}
.status-left {
  display: flex;
  align-items: center;
  gap: 12px;
}
.status-icon {
  font-size: 32px;
}
.status-icon.active {
  color: #67c23a;
}
.status-icon.paused {
  color: #e6a23c;
}
.status-title {
  font-size: 16px;
  font-weight: 600;
}
.status-sub {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
.section {
  margin-bottom: 16px;
}
.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}
.section-hint {
  font-size: 12px;
  color: #909399;
  font-weight: normal;
}
.qa-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 12px;
}
.qa-card {
  border-left: 4px solid #409eff;
}
.qa-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-weight: 600;
}
.add-card-btn {
  width: 100%;
  border: 1px dashed #409eff;
}
.button-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 12px;
}
.button-card {
  border-left: 4px solid #67c23a;
}
.footer-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  padding: 24px 0;
}
.cover-uploader :deep(.el-upload) {
  border: 1px dashed #d9d9d9;
  border-radius: 8px;
  cursor: pointer;
  overflow: hidden;
  width: 120px;
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: border-color 0.2s;
}
.cover-uploader :deep(.el-upload:hover) {
  border-color: #409eff;
}
.cover-preview {
  width: 120px;
  height: 120px;
  object-fit: cover;
  display: block;
}
.cover-uploader-icon {
  font-size: 28px;
  color: #8c939d;
}
.cover-tip {
  margin-top: 4px;
  font-size: 11px;
  color: #909399;
  text-align: center;
}
</style>
