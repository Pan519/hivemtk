<template>
  <div class="email-editor">
    <el-row :gutter="16">
      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.DragEditor.5.4b00ca') }}</span>
              <div>
                <el-button @click="previewDesktop">{{ $t('auto.DragEditor.6.745bb9') }}</el-button>
                <el-button @click="previewMobile">{{ $t('auto.DragEditor.7.e79ace') }}</el-button>
                <el-button type="primary" @click="save">{{ $t('auto.DragEditor.8.be5fbb') }}</el-button>
              </div>
            </div>
          </template>
          <div class="block-palette">
            <div
              v-for="block in blocks"
              :key="block.type"
              class="block-item"
              role="button"
              tabindex="0"
              :aria-label="`添加${block.label}块`"
              draggable="true"
              @dragstart="onDragStart($event, block)"
              @click="addBlock(block.type)"
              @keydown.enter.prevent="addBlock(block.type)"
              @keydown.space.prevent="addBlock(block.type)"
            >
              <el-icon><component :is="block.icon" /></el-icon>
              <span>{{ block.label }}</span>
            </div>
          </div>
          <div class="canvas" role="presentation" @drop="onDrop" @dragover.prevent>
            <div
              v-for="(block, i) in blocks_state"
              :key="block.id"
              class="canvas-block"
            >
              <component :is="resolveBlock(block.type)" v-bind="block" :editable="true" />
              <el-button-group size="small" class="block-actions">
                <el-button @click="moveUp(i)" :disabled="i === 0">↑</el-button>
                <el-button @click="moveDown(i)" :disabled="i === blocks_state.length - 1">↓</el-button>
                <el-button type="danger" @click="removeBlock(i)">{{ $t('auto.DragEditor.9.fe5c38') }}</el-button>
              </el-button-group>
            </div>
            <div v-if="blocks_state.length === 0" class="empty">{{ $t('auto.DragEditor.10.34c361') }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header><span>{{ $t('auto.DragEditor.11.bc3480') }}</span></template>
          <el-form label-width="80px">
            <el-form-item :label="$t('auto.DragEditor.1.9970ad')">
              <el-input v-model="subject" />
            </el-form-item>
            <el-form-item :label="$t('auto.DragEditor.2.88c952')">
              <el-input v-model="fromName" />
            </el-form-item>
            <el-form-item :label="$t('auto.DragEditor.3.b3e170')">
              <el-tag v-for="(t, i) in testEmails" :key="i" closable @close="testEmails.splice(i, 1)">{{ t }}</el-tag>
              <el-input v-model="newTestEmail" size="small" @keyup.enter="addTestEmail" placeholder="test@x.com" style="width: 200px" />
            </el-form-item>
            <el-form-item>
              <el-button @click="sendTest" :loading="sending">{{ $t('auto.DragEditor.12.99b3a7') }}</el-button>
            </el-form-item>
          </el-form>
          <el-divider />
          <h4>{{ $t('auto.DragEditor.13.4e1d5d') }}</h4>
          <el-tag v-for="v in builtinVars" :key="v.key" size="small" style="margin: 2px">
            {{ v.key }}
          </el-tag>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="previewVisible" :title="$t('auto.DragEditor.4.645dbc')" width="640px">
      <div class="preview-frame" v-html="compiledHTML" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n';
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Document, Picture, EditPen, PriceTag, Share, ChatLineRound, Loading
} from '@element-plus/icons-vue'
import mjml2html from 'mjml-browser'
import { http } from '@/utils/request'

const BLOCK_TYPES = [
  { type: 'text', label: t('auto.DragEditor.14.97d076'), icon: 'EditPen', default: { content: t('auto.DragEditor.15.657b63') } },
  { type: 'image', label: t('auto.DragEditor.16.20def7'), icon: 'Picture', default: { src: '', alt: '' } },
  { type: 'button', label: t('auto.DragEditor.17.fa9663'), icon: 'PriceTag', default: { text: t('auto.DragEditor.18.65e6e6'), url: '#', color: '#4F46E5' } },
  { type: 'divider', label: t('auto.DragEditor.19.8c47ec'), icon: 'Minus' },
  { type: 'social', label: t('auto.DragEditor.20.8b8588'), icon: 'Share', default: { networks: ['wechat', 'weibo'] } },
  { type: 'coupon', label: t('auto.DragEditor.21.2f3635'), icon: 'PriceTag', default: { code: 'SAVE10', amount: 10, expires: '30d' } }
]

const blocks = BLOCK_TYPES.map((b) => ({ ...b, icon: resolveIcon(b.icon) }))

const { t } = useI18n()

const blocks_state = ref([])
const subject = ref('')
const fromName = ref(t('auto.DragEditor.22.e6976d'))
const testEmails = ref([])
const newTestEmail = ref('')
const sending = ref(false)
const previewVisible = ref(false)
const compiledHTML = ref('')
const saving = ref(false)

const builtinVars = [
  { key: '{{customer.name}}' },
  { key: '{{order.id}}' },
  { key: '{{order.amount}}' },
  { key: '{{product.title}}' },
  { key: '{{agent.name}}' },
  { key: '{{coupon.code}}' }
]

function resolveIcon(name) {
  return { EditPen, Picture, PriceTag, Share, Document, ChatLineRound, Loading }[name] || Document
}

let _blockId = 1
function onDragStart(e, block) {
  e.dataTransfer.setData('blockType', block.type)
  e.dataTransfer.effectAllowed = 'copy'
}
function onDrop(e) {
  const type = e.dataTransfer.getData('blockType')
  if (!type) return
  addBlock(type)
}

// 键盘 / 点击的等价入口：拖放之外的可达替代（a11y）
function addBlock(type) {
  const tpl = BLOCK_TYPES.find((b) => b.type === type)
  if (!tpl) return
  blocks_state.value.push({
    id: _blockId++,
    type,
    ...(tpl?.default || {})
  })
}

function resolveBlock(type) {
  return {
    text: 'TextBlock',
    image: 'ImageBlock',
    button: 'ButtonBlock',
    divider: 'DividerBlock',
    social: 'SocialBlock',
    coupon: 'CouponBlock'
  }[type] || 'TextBlock'
}

function moveUp(i) { if (i > 0) blocks_state.value.splice(i - 1, 0, blocks_state.value.splice(i, 1)[0]) }
function moveDown(i) { if (i < blocks_state.value.length - 1) blocks_state.value.splice(i + 1, 0, blocks_state.value.splice(i, 1)[0]) }
function removeBlock(i) { blocks_state.value.splice(i, 1) }

function addTestEmail() {
  if (newTestEmail.value && /@/.test(newTestEmail.value)) {
    testEmails.value.push(newTestEmail.value)
    newTestEmail.value = ''
  }
}

async function sendTest() {
  if (testEmails.value.length === 0) return ElMessage.warning(t('auto.DragEditor.23.cf001e'))
  sending.value = true
  try {
    await http.post('/api/email/test-send', {
      subject: subject.value,
      html: compiledHTML.value,
      to: testEmails.value
    })
    ElMessage.success(t('auto.DragEditor.24.17e557'))
  } finally {
    sending.value = false
  }
}

async function save() {
  compiledHTML.value = compileToHTML()
  saving.value = true
  try {
    await http.post('/api/email/drafts', {
      subject: subject.value || t('auto.DragEditor.25.93db44'),
      content: compiledHTML.value
    })
    ElMessage.success(t('auto.DragEditor.26.22b831'))
  } catch (e) {
    ElMessage.error(t('auto.DragEditor.27.78d235') + (e?.message || e))
  } finally {
    saving.value = false
  }
}

function previewDesktop() {
  compiledHTML.value = compileToHTML()
  previewVisible.value = true
}
function previewMobile() {
  compiledHTML.value = '<style>body{max-width:375px;margin:0 auto;}</style>' + compileToHTML()
  previewVisible.value = true
}

// 邮件模板可能来自资产市场导入(他人设计),所有插值必须转义,
// URL 属性走协议白名单,否则预览/发送构成存储型 XSS 与属性注入
const esc = (s) => String(s ?? '')
  .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
const escUrl = (u) => {
  const s = esc(u)
  return /^(https?:|mailto:|tel:|#|\/)/i.test(s.trim()) ? s : '#'
}

function compileToHTML() {
  const parts = blocks_state.value.map((b) => {
    if (b.type === 'text') return `<mj-text>${esc(b.content)}</mj-text>`
    if (b.type === 'image') return `<mj-image src="${escUrl(b.src)}" alt="${esc(b.alt)}" />`
    if (b.type === 'button') return `<mj-button href="${escUrl(b.url)}" background-color="${esc(b.color || '#4F46E5')}">${esc(b.text || 'Button')}</mj-button>`
    if (b.type === 'divider') return '<mj-divider />'
    if (b.type === 'social') return `<mj-social><mj-social-element name="${esc(b.networks?.[0] || 'wechat')}" /></mj-social>`
    if (b.type === 'coupon') return `<mj-text>优惠券：${esc(b.code)}（减 ${esc(b.amount)}，${esc(b.expires)} 内有效）</mj-text>`
    return ''
  });
  const mjml = `<mjml><mj-body>${parts.join('\n')}</mj-body></mjml>`
  try {
    return mjml2html(mjml).html
  } catch (e) {
    return `<pre>${esc(e.message)}</pre>`
  }
}
</script>

<style scoped>
.email-editor { padding: 16px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.block-palette {
  display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 16px;
  padding: 12px; background: #F8FAFC; border-radius: 6px;
}
.block-item {
  display: flex; align-items: center; gap: 4px;
  padding: 6px 12px; background: #fff;
  border: 1px solid #E2E8F0; border-radius: 4px;
  cursor: grab; font-size: 12px;
}
.canvas {
  min-height: 400px; padding: 16px;
  background: #fff; border: 2px dashed #E2E8F0; border-radius: 6px;
}
.canvas-block { position: relative; margin: 8px 0; padding: 12px; border: 1px solid #E2E8F0; border-radius: 4px; }
.canvas-block:hover { border-color: #6366F1; }
.block-actions { position: absolute; right: 8px; top: 8px; }
.empty { text-align: center; color: #94A3B8; padding: 60px 0; }
.preview-frame { border: 1px solid #E2E8F0; padding: 16px; max-height: 600px; overflow: auto; }
</style>
