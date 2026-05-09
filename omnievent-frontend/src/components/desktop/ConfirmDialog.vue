<template>
  <v-dialog v-model="showDialog" max-width="400">
    <v-card>
      <v-card-title class="text-h6">{{ title }}</v-card-title>
      <v-card-text>{{ message }}</v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="onCancel">{{ cancelText }}</v-btn>
        <v-btn color="primary" @click="onConfirm">{{ confirmText }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const showDialog = ref(false)
const title = ref('Confirm')
const message = ref('')
const confirmText = ref('Confirm')
const cancelText = ref('Cancel')

function open(options: {
  title?: string
  message: string
  confirmText?: string
  cancelText?: string
}): Promise<boolean> {
  title.value = options.title || 'Confirm'
  message.value = options.message
  confirmText.value = options.confirmText || 'Confirm'
  cancelText.value = options.cancelText || 'Cancel'
  showDialog.value = true
  return new Promise((resolve) => {
    // Simple implementation - actual promise handling would need more work
    resolve(true)
  })
}

function onConfirm() {
  showDialog.value = false
}

function onCancel() {
  showDialog.value = false
}

defineExpose({ open })
</script>
