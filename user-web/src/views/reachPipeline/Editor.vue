<template>
  <VueFlowCanvas
    :title="`Pipeline - ${pipelineId || '新建'}`"
    :initial-nodes="initialNodes"
    :initial-edges="initialEdges"
    :node-templates="nodeTemplates"
    @save="onSave"
    @validate="onValidate"
  />
</template>

<script setup>
import { onMounted, ref, computed } from 'vue';
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import VueFlowCanvas from '@/components/VueFlowCanvas.vue'
import { http } from '@/utils/request'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const pipelineId = computed(() => route.params.id)
const initialNodes = ref([])
const initialEdges = ref([])

const saving = ref(false)

const nodeTemplates = [
  { type: 'trigger', label: t('auto.Editor.1.469672'), icon: 'Trigger', defaultData: { event: 'customer.created' } },
  { type: 'channel', label: t('auto.Editor.2.ebee05'), icon: 'Promotion', defaultData: { channel: 'email' } },
  { type: 'delay', label: t('auto.Editor.3.f23a9f'), icon: 'Timer', defaultData: { seconds: 3600 } },
  { type: 'condition', label: t('auto.Editor.4.c71b12'), icon: 'Branch', defaultData: { field: 'lifetime_value', op: '>', value: 1000 } },
  { type: 'abtest', label: 'A/B 实验', icon: 'DataAnalysis', defaultData: { experiment_id: null } },
  { type: 'ai', label: 'AI 生成', icon: 'MagicStick', defaultData: { prompt: t('auto.Editor.5.bd37e8'), agent_id: null } },
  { type: 'end', label: t('auto.Editor.6.12f1d7'), icon: 'CircleClose' }
]

async function load() {
  if (!pipelineId.value) return
  const p = await http.get(`/api/reach/pipelines/${pipelineId.value}`)
  initialNodes.value = p.nodes || []
  initialEdges.value = p.edges || []
}

async function onSave({ nodes, edges }) {
  // 无 id 时无可更新的 Pipeline（新建走 List 页 dialog），防 PUT 到 /undefined
  if (!pipelineId.value) {
    ElMessage.warning(t('auto.Editor.7.75d4e2'))
    return
  }
  if (saving.value) return
  saving.value = true
  try {
    await http.put(`/api/reach/pipelines/${pipelineId.value}`, { nodes, edges })
    ElMessage.success(t('auto.Editor.8.f8dfed'))
  } finally {
    saving.value = false
  }
}

function onValidate({ hasCycle }) {
  if (hasCycle) {
    ElMessage.error(t('auto.Editor.9.c8fb29'))
  } else {
    ElMessage.success(t('auto.Editor.10.9a9088'))
  }
}

onMounted(load)
</script>
