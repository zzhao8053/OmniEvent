<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card class="elevation-12">
          <v-card-title class="text-center">
            <h2>忘记密码</h2>
          </v-card-title>
          <v-card-text>
            <p class="text-body-2 text-center mb-4">
              请输入您的注册邮箱，我们将发送一封重置密码链接到您的邮箱
            </p>
            <v-form ref="formRef" @submit.prevent="handleRequestReset">
              <v-text-field
                v-model="email"
                label="邮箱"
                :rules="[rules.required, rules.email]"
                variant="outlined"
                prepend-inner-icon="mdi-email"
                autocomplete="email"
                :disabled="loading"
              />
              <v-alert v-if="successMessage" type="success" class="mb-4">
                {{ successMessage }}
              </v-alert>
              <v-alert v-if="errorMessage" type="error" class="mb-4">
                {{ errorMessage }}
              </v-alert>
              <v-btn
                type="submit"
                color="primary"
                block
                :loading="loading"
                size="large"
                :disabled="!email"
              >
                发送重置链接
              </v-btn>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-spacer />
            <v-btn variant="text" @click="$router.push('/login')">
              返回登录
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()

const email = ref('')
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const formRef = ref<any>(null)

const rules = {
  required: (v: string) => !!v || '此字段为必填项',
  email: (v: string) => /.+@.+\..+/.test(v) || '请输入有效的邮箱地址'
}

async function handleRequestReset() {
  const { valid } = await formRef.value.validate()
  if (!valid) return

  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    await axios.post('/api/forget_password/request.json', {
      email: email.value
    })
    successMessage.value = '重置密码链接已发送到您的邮箱'
  } catch (error: any) {
    errorMessage.value = error.response?.data?.errorMessage || '发送失败，请稍后重试'
  } finally {
    loading.value = false
  }
}
</script>
