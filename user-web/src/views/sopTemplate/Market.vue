<template>
  <div class="sop-market-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ $t('auto.Market.0a638b21') }}</h2>
          <p class="subtitle">{{ $t('auto.Market.ab0f3e86') }}</p>
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
        <el-form-item :label="$t('auto.Market.9699a50e')">
          <el-input
            v-model="filter.keyword"
            :placeholder="$t('auto.Market.ca75e9db')"
            clearable
            style="width: 220px"
            @keyup.enter="onSearch"
            @clear="onSearch"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.Market.d54920fa')">
          <el-select
            v-model="filter.category"
            :placeholder="$t('auto.Market.4f9a050f')"
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
          <el-button @click="resetFilter">{{ $t('auto.Market.4b9c3271') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    
    <div class="template-grid" v-loading="loading">
      <el-empty
        v-if="!loading && filteredTemplates.length === 0"
        :description="$t('auto.Market.5c9a5577')"
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
          <span class="scenario-label">{{ $t('auto.Market.fa272f67') }}</span>
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
            <span class="stat-label">{{ $t('auto.Market.764aca4e') }}</span>
            <el-tag size="small" type="warning">{{ tpl.intent }}</el-tag>
          </div>
          <div class="stat-item">
            <span class="stat-label">{{ $t('auto.Market.d316330b') }}</span>
            <el-tag size="small" type="info">{{ tpl.stage }}</el-tag>
          </div>
        </div>

        <div class="tpl-preview">
          <div class="preview-label">{{ $t('auto.Market.b638e662') }}</div>
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
      ::title="$t('auto.Market.013c0e71')"
      width="720px"
      :close-on-click-modal="false"
    >
      <template v-if="currentTpl">
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="$t('auto.Market.9118f5d2')">{{ currentTpl.code }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Market.d54920fa')">
            {{ getCategoryLabel(currentTpl.category) }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Market.269f953b')">{{ currentTpl.intent }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Market.73c7662e')">{{ currentTpl.stage }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Market.ee8ecb9e')">{{ currentTpl.priority }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Market.23d4eb1b')">{{ (currentTpl.confidence * 100).toFixed(0) }}%</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Market.3bdd08ad')" :span="2">
            {{ currentTpl.description }}
          </el-descriptions-item>
        </el-descriptions>

        <h4 style="margin-top: 16px">{{ $t('auto.Market.7bcd316c') }}</h4>
        <pre class="full-template">{{ currentTpl.template }}</pre>

        <h4 style="margin-top: 16px">{{ $t('auto.Market.ebbcc1d2') }}</h4>
        <ul class="vars-list">
          <li v-for="(v, idx) in currentTpl.vars" :key="idx">
            <code v-text="'{{' + v.name + '}}'"></code> - {{ v.desc }}
          </li>
        </ul>
      </template>

      <template #footer>
        <el-button @click="detailVisible = false">{{ $t('auto.Market.b15d9127') }}</el-button>
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
    name: t('auto.Market.82939b77'),
    description: t('auto.Market.14150688'),
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
    scenarios: [t('auto.Market.a9001500'), t('auto.Market.725965cf'), t('auto.Market.2d999029')],
    template: `您好 {{customer_name}}，感谢您对{{product_name}}的关注！
我是{{agent_name}}，专注于为{{industry}}客户提供解决方案。

针对您刚才咨询的{{topic}}，我整理了以下资料：
1. 产品功能介绍文档
2. 同行业客户案例（{{case_count}}个）
3. 试用账号（有效期7天）

请问您方便本周{{time_slot}}安排一次 15 分钟的线上沟通吗？`,
    vars: [
      { name: 'customer_name', desc: t('auto.Market.d1bdc43e') },
      { name: 'product_name', desc: t('auto.Market.ff685527') },
      { name: 'agent_name', desc: t('auto.Market.907f275f') },
      { name: 'industry', desc: t('auto.Market.71d93c20') },
      { name: 'topic', desc: t('auto.Market.faf1939e') },
      { name: 'case_count', desc: t('auto.Market.bfe25887') },
      { name: 'time_slot', desc: t('auto.Market.c2b7ca1d') }
    ]
  },
  {
    code: 'aftersale_service',
    name: t('auto.Market.3508dec9'),
    description: t('auto.Market.aa747f7d'),
    category: 'aftersale',
    icon: markRaw(Service),
    color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    intent: 'support',
    stage: 'middle',
    priority: 90,
    confidence: 0.9,
    imported: false,
    _importing: false,
    scenarios: [t('auto.Market.4cd514f9'), t('auto.Market.296a3154'), t('auto.Market.e80c6e86')],
    template: `{{customer_name}}您好，我是售后客服{{agent_name}}。

您反馈的{{issue_type}}问题，工单号 #{{ticket_id}}，我们已收到。
当前处理进度：
- 处理人：{{handler}}
- 预计解决时间：{{eta}}
- 补偿方案：{{compensation}}

如有任何疑问，请直接回复本消息，我会第一时间跟进。
您也可以拨打 400-xxx-xxxx 转工单号查询进度。`,
    vars: [
      { name: 'customer_name', desc: t('auto.Market.d1bdc43e') },
      { name: 'agent_name', desc: t('auto.Market.907f275f') },
      { name: 'issue_type', desc: t('auto.Market.d5a57c1c') },
      { name: 'ticket_id', desc: t('auto.Market.769eb37b') },
      { name: 'handler', desc: t('auto.Market.4c8c9d4f') },
      { name: 'eta', desc: t('auto.Market.5d4fd27b') },
      { name: 'compensation', desc: t('auto.Market.ca4f8f5d') }
    ]
  },
  {
    code: 'festival_marketing',
    name: t('auto.Market.54429779'),
    description: t('auto.Market.a007b75b'),
    category: 'marketing',
    icon: markRaw(Present),
    color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
    intent: 'purchase',
    stage: 'initial',
    priority: 70,
    confidence: 0.75,
    imported: false,
    _importing: false,
    scenarios: [t('auto.Market.506fb97d'), t('auto.Market.d920c676'), t('auto.Market.92a663d9')],
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
      { name: 'festival_name', desc: t('auto.Market.a0752060') },
      { name: 'customer_name', desc: t('auto.Market.d1bdc43e') },
      { name: 'threshold', desc: t('auto.Market.13d6dffc') },
      { name: 'discount', desc: t('auto.Market.2fc887d1') },
      { name: 'coupon_code', desc: t('auto.Market.390c08e7') },
      { name: 'deadline', desc: t('auto.Market.468c5447') }
    ]
  },
  {
    code: 'holiday_greeting',
    name: t('auto.Market.ffc9c1d2'),
    description: t('auto.Market.2083e7b0'),
    category: 'marketing',
    icon: markRaw(Sunny),
    color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    intent: 'greeting',
    stage: 'initial',
    priority: 50,
    confidence: 0.7,
    imported: false,
    _importing: false,
    scenarios: [t('auto.Market.1e61a8c7'), t('auto.Market.81f04166'), t('auto.Market.6da69cab')],
    template: `🌕 {{customer_name}}，{{holiday_name}}快乐！

值此{{holiday_name}}佳节，{{company_name}}全体员工
向您和家人致以最诚挚的祝福 🎊

【节日小贴士】
{{tip}}

我们将继续为您提供更优质的服务，
期待与您在新的一年继续同行！

——{{company_name}}`,
    vars: [
      { name: 'customer_name', desc: t('auto.Market.d1bdc43e') },
      { name: 'holiday_name', desc: t('auto.Market.a0752060') },
      { name: 'company_name', desc: t('auto.Market.079627c1') },
      { name: 'tip', desc: t('auto.Market.1fe1bd24') }
    ]
  },
  {
    code: 'new_customer_conversion',
    name: t('auto.Market.da4d371c'),
    description: t('auto.Market.20bd1c40'),
    category: 'presale',
    icon: markRaw(UserFilled),
    color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    intent: 'follow_up',
    stage: 'late',
    priority: 85,
    confidence: 0.8,
    imported: false,
    _importing: false,
    scenarios: [t('auto.Market.78e92e80'), t('auto.Market.d8eca889'), t('auto.Market.8500cb80')],
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
      { name: 'customer_name', desc: t('auto.Market.d1bdc43e') },
      { name: 'product_name', desc: t('auto.Market.c6afa029') },
      { name: 'step', desc: t('auto.Market.01a85850') },
      { name: 'progress_bar', desc: t('auto.Market.e2e33e2c') },
      { name: 'next_step', desc: t('auto.Market.95d760ff') },
      { name: 'discount', desc: t('auto.Market.65e9d3aa') },
      { name: 'next_action', desc: t('auto.Market.38e09d96') },
      { name: 'agent_name', desc: t('auto.Market.907f275f') }
    ]
  }
]);

const loading = ref(false);
const filter = ref({ keyword: '', category: '' })
const detailVisible = ref(false)
const currentTpl = ref(null)

const categories = [
  { value: 'presale', label: t('auto.Market.63d2ec80') },
  { value: 'aftersale', label: t('auto.Market.59bd6897') },
  { value: 'marketing', label: t('auto.Market.916801e9') }
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
    console.warn(t('auto.Market.e5122446'), e)
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
    ElMessage.error(t('auto.Market.013db1b1') + (e?.message || t('auto.Market.974e7484')))
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
