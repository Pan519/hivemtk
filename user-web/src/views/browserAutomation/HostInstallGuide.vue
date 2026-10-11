<template>
  <div class="guide">
    <p class="lead">{{ $t('auto.HostInstallGuide.1.973ffb') }}</p>
    <ol class="steps">
      <li>{{ $t('auto.HostInstallGuide.2.03f04f') }}<code>cd user-web/browser_automation && npm install && npm run build</code>
        <div class="note">{{ $t('auto.HostInstallGuide.3.2b83f6') }}<code>dist</code>{{ $t('auto.HostInstallGuide.4.11ae6e') }}</div></li>
      <li>chrome://extensions → 打开「开发者模式」→「加载已解压的扩展程序」→ 选择
        <code>user-web/browser_automation/dist</code></li>
      <li>{{ $t('auto.HostInstallGuide.5.51105c') }}<code>user-server/cmd/nm-host/install.sh</code>
        <div class="note">{{ $t('auto.HostInstallGuide.6.d009e5') }}<code>key</code>{{ $t('auto.HostInstallGuide.7.ab4aa9') }}</div></li>
      <li>{{ $t('auto.HostInstallGuide.8.bb57d8') }}<code>~/.hivemtk/nm_host.conf</code>{{ $t('auto.HostInstallGuide.9.01d7aa') }}<code>token=</code> 行
        <div class="note">{{ $t('auto.HostInstallGuide.10.983005') }}</div></li>
    </ol>
    <p class="lead">{{ $t('auto.HostInstallGuide.11.cfe7bf') }}</p>
  </div>
</template>

<style scoped>
.guide { line-height: 1.8; }
.lead { margin: 0 0 6px; }
.steps { margin: 0 0 6px 18px; padding: 0; }
.steps li { margin-bottom: 6px; }
.note { color: #909399; font-size: 12px; line-height: 1.6; }
</style>
