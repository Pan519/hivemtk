<template>
  <div class="sop-market-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ $t('auto.Market.13.0a638b') }}</h2>
          <p class="subtitle">{{ $t('auto.Market.14.ab0f3e') }}</p>
        </div>
        <div class="header-stats">
          <el-tag type="info" size="large">
            <el-icon><Collection /></el-icon>
            共 {{ builtInTemplates.length }} 个内置模板
          </el-tag>
        </div>
      </div>
    </el-card>

    
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filter" @submit.prevent>
        <el-form-item :label="$t('auto.Market.1.9699a5')">
          <el-input
            v-model="filter.keyword"
            :placeholder="$t('auto.Market.2.ca75e9')"
            clearable
            style="width: 220px"
            @keyup.enter="onSearch"
            @clear="onSearch"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.Market.3.d54920')">
          <el-select
            v-model="filter.category"
            :placeholder="$t('auto.Market.4.4f9a05')"
            clearable
            style="width: 180px"
            @change="onSearch"
          >
            <el-option
              v-for="c in categories"
              :key="c.value"
              :label="c.label"
              :value="c.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetFilter">{{ $t('auto.Market.15.4b9c32') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    
    <div class="template-grid" v-loading="loading">
      <el-empty
        v-if="!loading && filteredTemplates.length === 0"
        :description="$t('auto.Market.5.5c9a55')"
      />
      <el-card
        v-for="tpl in filteredTemplates"
        :key="tpl.code"
        class="tpl-card"
        shadow="hover"
      >
        <div class="tpl-card-header">
          <div class="tpl-icon" :style="{ background: tpl.color }">
            <el-icon :size="24"><component :is="tpl.icon" /></el-icon>
          </div>
          <div class="tpl-title-area">
            <h3 class="tpl-title">{{ tpl.name }}</h3>
            <el-tag :type="getCategoryTagType(tpl.category)" size="small">
              {{ getCategoryLabel(tpl.category) }}
            </el-tag>
          </div>
        </div>

        <p class="tpl-desc">{{ tpl.description }}</p>

        <div class="tpl-scenarios">
          <span class="scenario-label">{{ $t('auto.Market.16.fa272f') }}</span>
          <el-tag
            v-for="sc in tpl.scenarios"
            :key="sc"
            type="info"
            size="small"
            effect="plain"
            class="scenario-tag"
          >
            {{ sc }}
          </el-tag>
        </div>

        <div class="tpl-stats">
          <div class="stat-item">
            <span class="stat-label">{{ $t('auto.Market.17.764aca') }}</span>
            <el-tag size="small" type="warning">{{ tpl.intent }}</el-tag>
          </div>
          <div class="stat-item">
            <span class="stat-label">{{ $t('auto.Market.18.d31633') }}</span>
            <el-tag size="small" type="info">{{ tpl.stage }}</el-tag>
          </div>
        </div>

        <div class="tpl-preview">
          <div class="preview-label">{{ $t('auto.Market.19.b638e6') }}</div>
          <pre class="preview-content">{{ truncate(tpl.template, 140) }}</pre>
        </div>

        <div class="tpl-actions">
          <el-button
            v-if="tpl.imported"
            type="success"
            disabled
            size="default"
            style="flex: 1"
          >
            <el-icon><Check /></el-icon>
            已导入
          </el-button>
          <el-button
            v-else
            type="primary"
            size="default"
            :loading="tpl._importing"
            style="flex: 1"
            @click="onImport(tpl)"
          >
            <el-icon><Download /></el-icon>
            一键导入
          </el-button>
          <el-button
            size="default"
            @click="showDetail(tpl)"
          >
            <el-icon><View /></el-icon>
            详情
          </el-button>
        </div>
      </el-card>
    </div>

    
    <el-dialog
      v-model="detailVisible"
      ::title="$t('auto.Market.6.013c0e')"
      width="720px"
      :close-on-click-modal="false"
    >
      <template v-if="currentTpl">
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="$t('auto.Market.7.9118f5')">{{ currentTpl.code }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Market.3.d54920')">
            {{ getCategoryLabel(currentTpl.category) }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Market.8.269f95')">{{ currentTpl.intent }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Market.9.73c766')">{{ currentTpl.stage }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Market.10.ee8ecb')">{{ currentTpl.priority }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Market.11.23d4eb')">{{ (currentTpl.confidence * 100).toFixed(0) }}%</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Market.12.3bdd08')" :span="2">
            {{ currentTpl.description }}
          </el-descriptions-item>
        </el-descriptions>

        <h4 style="margin-top: 16px">{{ $t('auto.Market.20.7bcd31') }}</h4>
        <pre class="full-template">{{ currentTpl.template }}</pre>

        <h4 style="margin-top: 16px">{{ $t('auto.Market.21.ebbcc1') }}</h4>
        <ul class="vars-list">
          <li v-for="(v, idx) in currentTpl.vars" :key="idx">
            <code v-text="'{{' + v.name + '}}'"></code> - {{ v.desc }}
          </li>
        </ul>
      </template>

      <template #footer>
        <el-button @click="detailVisible = false">{{ $t('auto.Market.22.b15d91') }}</el-button>
        <el-button
          v-if="currentTpl && !currentTpl.imported"
          type="primary"
          :loading="currentTpl._importing"
          @click="onImport(currentTpl)"
        >
          一键导入
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, markRaw, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Collection,
  Search,
  View,
  Download,
  Check,
  Phone,
  Service,
  Present,
  Sunny,
  UserFilled,
  ChatLineRound
} from '@element-plus/icons-vue'
import { sopTemplateApi } from '@/api/sopTemplate'

const { t } = useI18n()

const builtInTemplates = ref([
  {
    code: 'pre_sale_followup',
    name: t('auto.Market.23.82939b'),
    description: t('auto.Market.24.141506'),
    category: 'presale',
    // icon 是 Vue 组件对象。builtInTemplates 是深响应 ref（下面还要就地写 t.imported，
    // 不能改用 shallowRef，否则模板 :116 的 v-if="tpl.imported" 会失去响应性），
    // 组件对象落进去会被 reactive 代理，模板 <component :is="tpl.icon" /> 取到代理组件
    // 就会触发 "Vue received a Component that was made a reactive object"（实测 109 条）。
    // markRaw 与 dashboardScreen/List.vue 同一修法。
    icon: markRaw(Phone),
    color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    intent: 'inquiry',
    stage: 'initial',
    priority: 80,
    confidence: 0.85,
    imported: false,
    _importing: false,
    scenarios: [t('auto.Market.25.a90015'), t('auto.Market.26.725965'), t('auto.Market.27.2d9990')],
    template: `您好 {{customer_name}}，感谢您对{{product_name}}的关注！
我是{{agent_name}}，专注于为{{industry}}客户提供解决方案。

针对您刚才咨询的{{topic}}，我整理了以下资料：
1. 产品功能介绍文档
2. 同行业客户案例（{{case_count}}个）
3. 试用账号（有效期7天）

请问您方便本周{{time_slot}}安排一次 15 分钟的线上沟通吗？`,
    vars: [
      { name: 'customer_name', desc: t('auto.Market.28.d1bdc4') },
      { name: 'product_name', desc: t('auto.Market.29.ff6855') },
      { name: 'agent_name', desc: t('auto.Market.30.907f27') },
      { name: 'industry', desc: t('auto.Market.31.71d93c') },
      { name: 'topic', desc: t('auto.Market.32.faf193') },
      { name: 'case_count', desc: t('auto.Market.33.bfe258') },
      { name: 'time_slot', desc: t('auto.Market.34.c2b7ca') }
    ]
  },
  {
    code: 'aftersale_service',
    name: t('auto.Market.35.3508de'),
    description: t('auto.Market.36.aa747f'),
    category: 'aftersale',
    icon: markRaw(Service),
    color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    intent: 'support',
    stage: 'middle',
    priority: 90,
    confidence: 0.9,
    imported: false,
    _importing: false,
    scenarios: [t('auto.Market.37.4cd514'), t('auto.Market.38.296a31'), t('auto.Market.39.e80c6e')],
    template: `{{customer_name}}您好，我是售后客服{{agent_name}}。

您反馈的{{issue_type}}问题，工单号 #{{ticket_id}}，我们已收到。
当前处理进度：
- 处理人：{{handler}}
- 预计解决时间：{{eta}}
- 补偿方案：{{compensation}}

如有任何疑问，请直接回复本消息，我会第一时间跟进。
您也可以拨打 400-xxx-xxxx 转工单号查询进度。`,
    vars: [
      { name: 'customer_name', desc: t('auto.Market.28.d1bdc4') },
      { name: 'agent_name', desc: t('auto.Market.30.907f27') },
      { name: 'issue_type', desc: t('auto.Market.40.d5a57c') },
      { name: 'ticket_id', desc: t('auto.Market.41.769eb3') },
      { name: 'handler', desc: t('auto.Market.42.4c8c9d') },
      { name: 'eta', desc: t('auto.Market.43.5d4fd2') },
      { name: 'compensation', desc: t('auto.Market.44.ca4f8f') }
    ]
  },
  {
    code: 'festival_marketing',
    name: t('auto.Market.45.544297'),
    description: t('auto.Market.46.a007b7'),
    category: 'marketing',
    icon: markRaw(Present),
    color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
    intent: 'purchase',
    stage: 'initial',
    priority: 70,
    confidence: 0.75,
    imported: false,
    _importing: false,
    scenarios: [t('auto.Market.47.506fb9'), t('auto.Market.48.d920c6'), t('auto.Market.49.92a663')],
    template: `🎉 {{festival_name}} 限时特惠！

亲爱的{{customer_name}}：
{{festival_name}}专属福利已到账，错过等一年 ✨

【尊享特权】
- 全场满 {{threshold}} 减 {{discount}}
- 爆款商品低至 5 折
- 专属客服 1v1 服务

👉 立即抢购：[点击领取优惠券 {{coupon_code}}]
⏰ 活动截止：{{deadline}}

祝您和家人{{festival_name}}快乐！`,
    vars: [
      { name: 'festival_name', desc: t('auto.Market.50.a07520') },
      { name: 'customer_name', desc: t('auto.Market.28.d1bdc4') },
      { name: 'threshold', desc: t('auto.Market.51.13d6df') },
      { name: 'discount', desc: t('auto.Market.52.2fc887') },
      { name: 'coupon_code', desc: t('auto.Market.53.390c08') },
      { name: 'deadline', desc: t('auto.Market.54.468c54') }
    ]
  },
  {
    code: 'holiday_greeting',
    name: t('auto.Market.55.ffc9c1'),
    description: t('auto.Market.56.2083e7'),
    category: 'marketing',
    icon: markRaw(Sunny),
    color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    intent: 'greeting',
    stage: 'initial',
    priority: 50,
    confidence: 0.7,
    imported: false,
    _importing: false,
    scenarios: [t('auto.Market.57.1e61a8'), t('auto.Market.58.81f041'), t('auto.Market.59.6da69c')],
    template: `🌕 {{customer_name}}，{{holiday_name}}快乐！

值此{{holiday_name}}佳节，{{company_name}}全体员工
向您和家人致以最诚挚的祝福 🎊

【节日小贴士】
{{tip}}

我们将继续为您提供更优质的服务，
期待与您在新的一年继续同行！

——{{company_name}}`,
    vars: [
      { name: 'customer_name', desc: t('auto.Market.28.d1bdc4') },
      { name: 'holiday_name', desc: t('auto.Market.50.a07520') },
      { name: 'company_name', desc: t('auto.Market.60.079627') },
      { name: 'tip', desc: t('auto.Market.61.1fe1bd') }
    ]
  },
  {
    code: 'new_customer_conversion',
    name: t('auto.Market.62.da4d37'),
    description: t('auto.Market.63.20bd1c'),
    category: 'presale',
    icon: markRaw(UserFilled),
    color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    intent: 'follow_up',
    stage: 'late',
    priority: 85,
    confidence: 0.8,
    imported: false,
    _importing: false,
    scenarios: [t('auto.Market.64.78e92e'), t('auto.Market.65.d8eca8'), t('auto.Market.66.8500cb')],
    template: `Hi {{customer_name}}，欢迎加入{{product_name}}！

我们注意到您已经完成了第{{step}}/4 步：
{{progress_bar}}

【下一步建议】{{next_step}}

📌 新客专属福利：
- 14 天专业版试用
- 1 对 1 培训（30 分钟）
- 首单立减 {{discount}} 元

立即体验：[点击开始 {{next_action}}]

有任何问题随时找我，我是{{agent_name}}。`,
    vars: [
      { name: 'customer_name', desc: t('auto.Market.28.d1bdc4') },
      { name: 'product_name', desc: t('auto.Market.67.c6afa0') },
      { name: 'step', desc: t('auto.Market.68.01a858') },
      { name: 'progress_bar', desc: t('auto.Market.69.e2e33e') },
      { name: 'next_step', desc: t('auto.Market.70.95d760') },
      { name: 'discount', desc: t('auto.Market.71.65e9d3') },
      { name: 'next_action', desc: t('auto.Market.72.38e09d') },
      { name: 'agent_name', desc: t('auto.Market.30.907f27') }
    ]
  }
]);

const loading = ref(false);
const filter = ref({ keyword: '', category: '' })
const detailVisible = ref(false)
const currentTpl = ref(null)

const categories = [
  { value: 'presale', label: t('auto.Market.73.63d2ec') },
  { value: 'aftersale', label: t('auto.Market.74.59bd68') },
  { value: 'marketing', label: t('auto.Market.75.916801') }
]

const filteredTemplates = computed(() => {
  let list = builtInTemplates.value
  if (filter.value.keyword) {
    const kw = filter.value.keyword.toLowerCase()
    list = list.filter(
      (t) =>
        t.name.toLowerCase().includes(kw) ||
        t.description.toLowerCase().includes(kw)
    )
  }
  if (filter.value.category) {
    list = list.filter((t) => t.category === filter.value.category)
  }
  return list
});

const truncate = (text, len) => {
  if (!text) return '-'
  return text.length > len ? text.slice(0, len) + '...' : text
};

const getCategoryLabel = (cat) => {
  const c = categories.find((x) => x.value === cat)
  return c ? c.label : cat
}

const getCategoryTagType = (cat) => {
  if (cat === 'presale') return 'warning'
  if (cat === 'aftersale') return 'danger'
  if (cat === 'marketing') return 'success'
  return 'info'
}

const loadImportedState = async () => {
  try {
    const res = await sopTemplateApi.list({ page: 1, page_size: 100 }).catch(() => null)
    const items = res?.list || []
    const importedCodes = new Set(
      items
        .filter((it) => it.code)
        .map((it) => it.code)
    )
    builtInTemplates.value.forEach((t) => {
      t.imported = importedCodes.has(t.code)
    })
  } catch (e) {
    console.warn(t('auto.Market.76.e51224'), e)
  }
};

const onSearch = () => {};

const resetFilter = () => {
  filter.value = { keyword: '', category: '' }
}

const showDetail = (tpl) => {
  currentTpl.value = tpl
  detailVisible.value = true
}

const onImport = async (tpl) => {
  tpl._importing = true
  try {
    await sopTemplateApi.create({
      code: tpl.code,
      name: tpl.name,
      intent: tpl.intent,
      stage: tpl.stage,
      template: tpl.template,
      vars: JSON.stringify(tpl.vars),
      priority: tpl.priority,
      confidence: tpl.confidence,
      enabled: true
    })
    tpl.imported = true
    ElMessage.success(`「${tpl.name}」已成功导入！`)
    if (detailVisible.value) {
      detailVisible.value = false
    }
  } catch (e) {
    ElMessage.error(t('auto.Market.77.013db1') + (e?.message || t('auto.Market.78.974e74')))
  } finally {
    tpl._importing = false
  }
}

onMounted(() => {
  loadImportedState()
})
</script>

<style scoped>
.sop-market-page {
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

.filter-card :deep(.el-card__body) {
  padding: 16px 20px 0 20px;
}

.template-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
  gap: 16px;
  padding: 4px;
}

.tpl-card {
  display: flex;
  flex-direction: column;
  transition: transform 0.2s;
}

.tpl-card:hover {
  transform: translateY(-2px);
}

.tpl-card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.tpl-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.tpl-title-area {
  flex: 1;
  min-width: 0;
}

.tpl-title {
  margin: 0 0 4px 0;
  font-size: 16px;
  font-weight: 600;
}

.tpl-desc {
  margin: 0 0 12px 0;
  color: var(--el-text-color-regular);
  font-size: 13px;
  line-height: 1.6;
  min-height: 40px;
}

.tpl-scenarios {
  margin-bottom: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  align-items: center;
}

.scenario-label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-right: 4px;
}

.scenario-tag {
  margin: 0;
}

.tpl-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
  padding: 8px 12px;
  background: var(--el-fill-color-light);
  border-radius: 4px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
}

.stat-label {
  color: var(--el-text-color-secondary);
}

.tpl-preview {
  margin-bottom: 12px;
  flex: 1;
}

.preview-label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-bottom: 4px;
}

.preview-content {
  margin: 0;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 4px;
  font-size: 12px;
  line-height: 1.6;
  color: #606266;
  font-family: 'SF Mono', Monaco, Consolas, monospace;
  max-height: 100px;
  overflow: hidden;
}

.tpl-actions {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.full-template {
  margin: 0;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 4px;
  font-size: 13px;
  line-height: 1.7;
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 320px;
  overflow-y: auto;
}

.vars-list {
  margin: 8px 0 0 0;
  padding-left: 20px;
  font-size: 13px;
  line-height: 1.8;
}

.vars-list code {
  background: #f0f9ff;
  color: #1890ff;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 12px;
}
</style>
