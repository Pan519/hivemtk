<template>
  <div class="faq-index-page">
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ $t('auto.Index.1.e571cd') }}</h2>
          <p class="subtitle">{{ $t('auto.Index.2.a4ea7b') }}</p>
        </div>
      </div>
    </el-card>

    <el-row :gutter="20">
      <el-col :span="8">
        <el-card shadow="hover" class="entry-card" @click="$router.push('/faq/list')">
          <div class="entry-icon"><el-icon><Notebook /></el-icon></div>
          <div class="entry-title">{{ $t('auto.Index.3.c7b1da') }}</div>
          <div class="entry-desc">{{ $t('auto.Index.4.5a1a28') }}</div>
          <el-button type="primary" plain>
            进入
            <el-icon class="el-icon--right"><ArrowRight /></el-icon>
          </el-button>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="entry-card">
          <div class="entry-icon"><el-icon><Aim /></el-icon></div>
          <div class="entry-title">{{ $t('auto.Index.5.c50cce') }}</div>
          <div class="entry-desc">{{ $t('auto.Index.6.9dc6c4') }}</div>
          <el-button type="primary" plain @click="$router.push('/faq/list')">
            打开
            <el-icon class="el-icon--right"><ArrowRight /></el-icon>
          </el-button>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="entry-card">
          <div class="entry-icon"><el-icon><Connection /></el-icon></div>
          <div class="entry-title">{{ $t('auto.Index.7.629702') }}</div>
          <div class="entry-desc">{{ $t('auto.Index.8.6067dc') }}</div>
          <el-button type="primary" plain @click="$router.push('/aiAgent/list')">
            去绑定
            <el-icon class="el-icon--right"><ArrowRight /></el-icon>
          </el-button>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { Notebook, Aim, Connection, ArrowRight } from '@element-plus/icons-vue'
</script>

<style scoped>
.faq-index-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.header-card :deep(.el-card__body) {
  padding: 18px 24px;
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

.entry-card {
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
  text-align: center;
  padding: 8px 0;
}

.entry-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.entry-icon {
  font-size: 36px;
  color: var(--el-color-primary);
  margin-bottom: 8px;
}

.entry-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 6px;
}

.entry-desc {
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.5;
  margin-bottom: 12px;
  min-height: 40px;
}
</style>
