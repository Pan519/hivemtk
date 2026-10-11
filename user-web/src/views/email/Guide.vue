<template>
  <div class="guide-container">
    <h1 class="guide-title">{{ $t('邮件营销使用指南') }}</h1>
    <div class="guide-intro">
      <p>{{ $t('auto.Guide.90485ce8') }}</p>
    </div>

    <el-steps :active="activeStep" finish-status="success" align-center class="steps-container">
      <el-step ::title="$t('auto.Guide.22d45263')" :description="$t('auto.Guide.b0ec481c')" />
      <el-step ::title="$t('auto.Guide.5c9b35b5')" :description="$t('auto.Guide.654ea942')" />
      <el-step ::title="$t('auto.Guide.2dc0a534')" :description="$t('auto.Guide.7706ff4a')" />
      <el-step ::title="$t('auto.Guide.a9bac747')" :description="$t('auto.Guide.16a61025')" />
    </el-steps>

    <el-row :gutter="20" class="step-cards">
      <el-col :span="6" v-for="(item, idx) in steps" :key="idx">
        <el-card shadow="hover" class="step-card" :class="{ active: activeStep === idx }" @click="activeStep = idx">
          <h3>{{ item.title }}</h3>
          <p class="desc">{{ item.description }}</p>
          <ul>
            <li v-for="(d, i) in item.details" :key="i">{{ d }}</li>
          </ul>
          <el-button type="primary" link @click.stop="goTo(item.path)">
            {{ $t('前往操作') }} <el-icon><ArrowRight /></el-icon>
          </el-button>
        </el-card>
      </el-col>
    </el-row>

    <el-alert ::title="$t('auto.Guide.d29c3db7')" type="info" :closable="false" show-icon class="faq">
      <p><strong>{{ $t('auto.Guide.2ae55257') }}</strong>{{ $t('auto.Guide.ff52b658') }}</p>
      <p><strong>{{ $t('auto.Guide.0ffcfb20') }}</strong> {{ $t('建议使用企业域名邮箱，避免营销话术，并合理设置发送频率。') }}</p>
      <p><strong>{{ $t('auto.Guide.8d4e899f') }}</strong> {{ $t('适当增加任务并发数，但需避免触发反垃圾邮件规则。') }}</p>
    </el-alert>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const router = useRouter()
const { t } = useI18n()

const activeStep = ref(0)

const steps = [
  {
    title: i18n.global.t('1. 配置 SMTP'),
    description: t('auto.Guide.5ce3b2a3'),
    details: [
      t('auto.Guide.fccb5b3b'),
      t('auto.Guide.2b193a89'),
      t('auto.Guide.7019e40c')
    ],
    path: '/email/smtp'
  },
  {
    title: i18n.global.t('2. 收件人列表'),
    description: t('auto.Guide.3ed4dd13'),
    details: [
      t('auto.Guide.ddcd34d4'),
      t('auto.Guide.f4dd3671'),
      t('auto.Guide.b8a79f22')
    ],
    path: '/email/list'
  },
  {
    title: i18n.global.t('3. 编写草稿'),
    description: t('auto.Guide.7706ff4a'),
    details: [
      t('auto.Guide.a28521ae'),
      t('auto.Guide.1b4f194c'),
      t('auto.Guide.81b9269a')
    ],
    path: '/email/drafts'
  },
  {
    title: i18n.global.t('4. 发送任务'),
    description: t('auto.Guide.3978ef5d'),
    details: [
      t('auto.Guide.8d0988e5'),
      t('auto.Guide.ef8d53db'),
      t('auto.Guide.2cfd19a6')
    ],
    path: '/email/jobs'
  }
]

const goTo = (path) => {
  router.push(path).catch(() => {})
}
</script>

<style scoped>
.guide-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}
.guide-title {
  text-align: center;
  margin-bottom: 20px;
  color: #333;
}
.guide-intro {
  background-color: #f0f7ff;
  padding: 15px 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  border: 1px solid #e1f0ff;
}
.steps-container {
  margin-bottom: 30px;
}
.step-cards {
  margin-bottom: 30px;
}
.step-card {
  cursor: pointer;
  transition: all 0.2s;
  height: 100%;
}
.step-card.active {
  border-color: #4F46E5;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.15);
}
.step-card h3 {
  margin-top: 0;
  color: #303133;
}
.step-card .desc {
  color: #606266;
  font-size: 13px;
  margin: 8px 0;
}
.step-card ul {
  padding-left: 18px;
  color: #606266;
  font-size: 13px;
  line-height: 1.8;
}
.faq {
  margin-top: 20px;
}
.faq p {
  margin: 4px 0;
  line-height: 1.6;
}
</style>
