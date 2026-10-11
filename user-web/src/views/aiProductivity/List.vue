<template>
  <div class="ai-productivity">
    <el-card class="page-card">
      <div class="page-header">
        <div>
          <h2 class="page-title">{{ $t('auto.List.9.063fc7') }}</h2>
          <p class="page-sub">{{ $t('auto.List.10.c43ffd') }}</p>
        </div>
        <el-button type="primary" @click="loadAll">{{ $t('auto.List.11.694fc5') }}</el-button>
      </div>

      <el-row :gutter="16" class="kpi-row">
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.1.7ed257')" :value="rep.total_conversations" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.2.cde02b')" :value="rep.ai_ratio" suffix="%" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.3.3871cd')" :value="Number(Number(rep.avg_response_time || 0).toFixed(1))" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.4.715054')" :value="rep.conversion_rate" suffix="%" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.5.9b925f')" :value="rep.ai_replies" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.6.58661d')" :value="rep.human_replies" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.7.2550ea')" :value="rep.total_conversions" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="LLM Token" :value="rep.llm_tokens" />
        </el-col>
      </el-row>

      <el-card shadow="never" :header="$t('auto.List.8.3c3478')" class="trend-card">
        <div ref="trendChart" class="chart-box"></div>
      </el-card>
    </el-card>
  </div>
</template>

<script>
import * as echarts from 'echarts'
import i18n from '@/i18n'
const t = i18n.global.t
import { safeInit } from '@/utils/echarts'
import AIProductivityApi from '@/api/aiProductivity'

export default {
  name: 'AIProductivityList',
  data() {
    return {
      rep: {},
      trend: [],
      trendChart: null,
      chartInst: null
    }
  },
  mounted() {
    this.chartInst = safeInit(this.$refs.trendChart)
    this.loadAll()
  },
  beforeUnmount() {
    if (this.chartInst) this.chartInst.dispose()
  },
  methods: {
    async loadAll() {
      const [rep, trendRes] = await Promise.all([
        AIProductivityApi.getReport({}),
        AIProductivityApi.getTrend({ days: 30 })
      ])
      this.rep = rep || {}
      const trendData = trendRes && Array.isArray(trendRes.trend) ? trendRes.trend : (Array.isArray(trendRes) ? trendRes : [])
      this.trend = t
      this.renderTrend()
    },
    renderTrend() {
      if (!this.chartInst) return
      const dates = this.trend.map((d) => d.date)
      const conv = this.trend.map((d) => d.conversations)
      const ai = this.trend.map((d) => d.ai_replies)
      const orders = this.trend.map((d) => d.conversions)
      this.chartInst.setOption({
        tooltip: { trigger: 'axis' },
        legend: { data: [t('auto.List.12.4d4b9c'), 'AI 回复数', t('auto.List.13.d4268d')] },
        xAxis: { type: 'category', data: dates },
        yAxis: { type: 'value' },
        series: [
          { name: t('auto.List.12.4d4b9c'), type: 'line', data: conv },
          { name: 'AI 回复数', type: 'line', data: ai },
          { name: t('auto.List.13.d4268d'), type: 'bar', data: orders }
        ]
      })
    }
  }
}
</script>

<style scoped>
.page-card { padding: 8px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.page-title { margin: 0; font-size: 20px; }
.page-sub { margin: 4px 0 0; color: #909399; font-size: 13px; }
.kpi-row { margin-bottom: 16px; }
.trend-card { margin-top: 8px; }
.chart-box { height: 340px; }
</style>
