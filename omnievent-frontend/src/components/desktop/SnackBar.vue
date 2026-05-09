<template>
  <v-snackbar
    v-model="show"
    :color="color"
    :timeout="timeout"
    location="top"
  >
    {{ message }}
    <template v-slot:actions>
      <v-btn variant="text" @click="show = false">
        {{ $t('close') }}
      </v-btn>
    </template>
  </v-snackbar>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const show = ref(false)
const message = ref('')
const color = ref('primary')
const timeout = ref(3000)

function open(msg: string, msgColor: string = 'primary', msgTimeout: number = 3000) {
  message.value = msg
  color.value = msgColor
  timeout.value = msgTimeout
  show.value = true
}

defineExpose({ open })
</script>
