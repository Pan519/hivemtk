<template>
  <div class="merchant-editor">
    
    <el-card class="status-bar" shadow="never">
      <div class="status-row">
        <div class="status-left">
          <el-icon v-if="form.status === 'active'" class="status-icon active"><Monitor /></el-icon>
          <el-icon v-else class="status-icon paused"><VideoPause /></el-icon>
          <div class="status-text">
            <div class="status-title">
              <span v-if="form.status === 'active'">{{ $t('auto.MerchantEditor.3dc4bdc4') }}</span>
              <span v-else>{{ $t('auto.MerchantEditor.e03a18e5') }}</span>
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
            <span>{{ $t('auto.MerchantEditor.0f8a0c04') }}</span>
            <span class="section-hint">{{ $t('auto.MerchantEditor.ee879b1c') }}</span>
          </div>
        </template>
        <el-row :gutter="20">
          <el-col :span="6">
            <el-form-item :label="$t('auto.MerchantEditor.641cfa0a')">
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
                >
                  <template #error>
                    <div class="img-error-fallback">{{ $t('图片加载失败') }}</div>
                  </template>
                </el-image>
                <el-icon v-else class="cover-uploader-icon"><Plus /></el-icon>
              </el-upload>
              <div class="cover-tip">{{ $t('auto.MerchantEditor.70baac73') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="18">
            <el-form-item :label="$t('auto.MerchantEditor.0acf6dc2')">
              <el-input v-model="form.title" :placeholder="$t('auto.MerchantEditor.5a5db2bd')" />
            </el-form-item>
            <el-form-item :label="$t('auto.MerchantEditor.69847d0a')">
              <el-input v-model="form.asset_id" :placeholder="$t('auto.MerchantEditor.01e0f1ef')" :disabled="isEdit" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('auto.MerchantEditor.0ad1ab6f')">
              <el-input v-model="form.campaign_name" :placeholder="$t('auto.MerchantEditor.a7a9c157')" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item :label="$t('auto.MerchantEditor.59e2db23')">
              <el-input v-model="form.discount_pct" :placeholder="$t('auto.MerchantEditor.8d811279')" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item :label="$t('auto.MerchantEditor.5623372c')">
              <el-input v-model="form.support_contact" :placeholder="$t('auto.MerchantEditor.2410033c')" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('auto.MerchantEditor.53732575')">
              <el-input v-model="form.shop_name" :placeholder="$t('auto.MerchantEditor.3c02ae66')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.MerchantEditor.be3d89df')">
              <el-input v-model="form.author" :placeholder="$t('auto.MerchantEditor.015d6ef6')" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      
      <el-card class="section" shadow="never">
        <template #header>
          <div class="section-header">
            <el-icon><Lock /></el-icon>
            <span>{{ $t('auto.MerchantEditor.1888faa8') }}</span>
            <span class="section-hint">{{ $t('auto.MerchantEditor.deac7199') }}</span>
          </div>
        </template>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item :label="$t('auto.MerchantEditor.c0fbfddc')">
              <el-select v-model="form.crisis_threshold" :placeholder="$t('auto.MerchantEditor.0a11d380')">
                <el-option :label="$t('auto.MerchantEditor.1b57b470')" value="2" />
                <el-option :label="$t('auto.MerchantEditor.1ff4bb24')" value="3" />
                <el-option :label="$t('auto.MerchantEditor.9e6086f1')" value="4" />
                <el-option :label="$t('auto.MerchantEditor.47db26ee')" value="5" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.MerchantEditor.c947d43e')">
              <el-select v-model="form.tone_level" :placeholder="$t('auto.MerchantEditor.8dfbf031')">
                <el-option :label="$t('auto.MerchantEditor.9a8e9954')" value="low" />
                <el-option :label="$t('auto.MerchantEditor.b51f755d')" value="medium" />
                <el-option :label="$t('auto.MerchantEditor.fd1d81f0')" value="high" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.MerchantEditor.9b4521a0')">
              <el-radio-group v-model="form.censorship_level">
                <el-radio label="strict">{{ $t('auto.MerchantEditor.15087bd7') }}</el-radio>
                <el-radio label="unlock">{{ $t('auto.MerchantEditor.50d7e2f7') }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('auto.MerchantEditor.e819d991')">
          <el-checkbox-group v-model="form.enabled_intents">
            <el-checkbox label="faq">{{ $t('auto.MerchantEditor.4bc10db3') }}</el-checkbox>
            <el-checkbox label="lead_capture">{{ $t('auto.MerchantEditor.535d4693') }}</el-checkbox>
            <el-checkbox label="human_transfer">{{ $t('auto.MerchantEditor.3e61a7bb') }}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-card>

      
      <el-card class="section" shadow="never">
        <template #header>
          <div class="section-header">
            <el-icon><ChatDotRound /></el-icon>
            <span>{{ $t('auto.MerchantEditor.08da94e8') }}</span>
            <span class="section-hint">{{ $t('auto.MerchantEditor.139adb0b') }}</span>
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
              <el-button type="danger" link @click="removeQACard(card.id)">{{ $t('auto.MerchantEditor.2f4aaddd') }}</el-button>
            </div>
            <el-form-item :label="$t('auto.MerchantEditor.ae37bc39')">
              <el-input v-model="card.trigger" :placeholder="$t('auto.MerchantEditor.02798fae')" />
            </el-form-item>
            <el-form-item :label="$t('auto.MerchantEditor.b4baaa14')">
              <el-input
                v-model="card.user_example"
                type="textarea"
                :rows="2"
                :placeholder="$t('auto.MerchantEditor.0b770e5d')"
              />
            </el-form-item>
            <el-form-item :label="$t('auto.MerchantEditor.9ac55b24')">
              <el-input
                v-model="card.reply"
                type="textarea"
                :rows="4"
                :placeholder="$t('auto.MerchantEditor.9c32be06')"
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
            <span>{{ $t('auto.MerchantEditor.4774fd33') }}</span>
            <span class="section-hint">{{ $t('auto.MerchantEditor.5d3663df') }}</span>
          </div>
        </template>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item :label="$t('auto.MerchantEditor.f60f4a0a')">
              <el-radio-group v-model="form.card_config.intent_type">
                <el-radio label="button_card">{{ $t('auto.MerchantEditor.aaa40b0b') }}</el-radio>
                <el-radio label="coupon">{{ $t('auto.MerchantEditor.c323c414') }}</el-radio>
                <el-radio label="handoff">{{ $t('auto.MerchantEditor.6586199c') }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="16">
            <el-form-item :label="$t('auto.MerchantEditor.83747b77')">
              <el-input v-model="form.card_config.product_image" :placeholder="$t('auto.MerchantEditor.2edb1436')" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-divider>{{ $t('auto.MerchantEditor.4a06d4c0') }}</el-divider>
        <div class="button-list">
          <el-card
            v-for="(btn, idx) in form.card_config.buttons"
            :key="idx"
            class="button-card"
            shadow="hover"
          >
            <el-row :gutter="10">
              <el-col :span="8">
                <el-form-item :label="$t('auto.MerchantEditor.ddefac28')">
                  <el-input v-model="btn.title" :placeholder="$t('auto.MerchantEditor.934285fe')" />
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item :label="$t('auto.MerchantEditor.db24c6ea')">
                  <el-select v-model="btn.action">
                    <el-option :label="$t('auto.MerchantEditor.e44a6327')" value="open_url" />
                    <el-option :label="$t('auto.MerchantEditor.9e7021aa')" value="call_api" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item v-if="btn.action === 'open_url'" :label="$t('auto.MerchantEditor.e44a6327')">
                  <el-input v-model="btn.url" placeholder="https://shopify.com/..." />
                </el-form-item>
                <el-form-item v-else :label="$t('auto.MerchantEditor.4b3dc025')">
                  <el-input v-model="btn.api_name" :placeholder="$t('auto.MerchantEditor.4a59e81d')" />
                </el-form-item>
              </el-col>
              <el-col :span="2">
                <el-button type="danger" link @click="removeButton(idx)">{{ $t('auto.MerchantEditor.2f4aaddd') }}</el-button>
              </el-col>
            </el-row>
            <el-form-item v-if="btn.action === 'call_api'" :label="$t('auto.MerchantEditor.8bc2cc75')">
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
        <el-button size="large" @click="$router.push('/asset-bundle/list')">{{ $t('auto.MerchantEditor.adcd1dd7') }}</el-button>
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
    ElMessage.error(t('auto.MerchantEditor.2cfc148d') + (e?.message || e))
  } finally {
    loading.value = false
  }
};

const handleCoverUpload = async ({ file }) => {
  try {
    const res = await uploadCover(file)
    form.cover_image = res.data?.url || res.url || res
    ElMessage.success(t('auto.MerchantEditor.e04a52d0'))
  } catch (e) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
};

const handleSave = async () => {
  if (!form.asset_id) {
    ElMessage.warning(t('auto.MerchantEditor.b9b9a97e'))
    return
  }
  if (!form.title) {
    ElMessage.warning(t('auto.MerchantEditor.ac8238fb'))
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
        ElMessage.success(t('auto.MerchantEditor.a7f38dfd'))
      } catch (e) {
        ElMessage.warning(t('auto.MerchantEditor.045c62db'))
      }
      if (!isEdit.value && data.asset_id) {
        router.replace(`/asset-bundle/merchant/${data.asset_id}`)
      }
    } else {
      ElMessage.success(t('auto.MerchantEditor.3b108349'))
    }
  } catch (e) {
    ElMessage.error(t('auto.MerchantEditor.40f90217') + (e?.message || e))
  } finally {
    saving.value = false
  }
};

const handleToggleStatus = async (action) => {
  if (!aid.value) {
    ElMessage.warning(t('auto.MerchantEditor.9558c46d'))
    return
  }
  try {
    const detailResp = await getBundleByAssetID(aid.value);
    const detail = detailResp?.data || detailResp
    if (!detail || !detail.id) {
      ElMessage.error(t('auto.MerchantEditor.23fe0dcb'))
      return
    }
    if (action === 'pause') {
      await disableBundle(detail.id)
      form.status = 'inactive'
      ElMessage.success(t('auto.MerchantEditor.76be2ffb'))
    } else {
      await enableBundle(detail.id)
      form.status = 'active'
      ElMessage.success(t('auto.MerchantEditor.74aa817e'))
    }
  } catch (e) {
    ElMessage.error(t('auto.MerchantEditor.1a4a3038') + (e?.message || e))
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
