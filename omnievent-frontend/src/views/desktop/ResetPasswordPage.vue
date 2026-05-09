<template>
  <div class="layout-wrapper">
    <router-link to="/">
      <div class="auth-logo d-flex align-start gap-x-3">
        <img alt="logo" class="login-page-logo" :src="APPLICATION_LOGO_PATH" />
        <h1 class="font-weight-medium leading-normal text-2xl">OmniEvent</h1>
      </div>
    </router-link>
    <v-row no-gutters class="auth-wrapper">
      <v-col cols="12" md="8" class="auth-image-background d-none d-md-flex align-center justify-center position-relative">
        <div class="d-flex auth-img-footer" v-if="!isDarkMode">
          <v-img class="img-with-direction" src="img/desktop/background.svg"/>
        </div>
        <div class="d-flex auth-img-footer" v-if="isDarkMode">
          <v-img class="img-with-direction" src="img/desktop/background-dark.svg"/>
        </div>
        <div class="d-flex align-center justify-center w-100 pt-10">
          <v-img class="img-with-direction" max-width="600px" src="img/desktop/people4.svg" v-if="!isDarkMode"/>
          <v-img class="img-with-direction" max-width="600px" src="img/desktop/people4-dark.svg" v-else-if="isDarkMode"/>
        </div>
      </v-col>
      <v-col cols="12" md="4" class="auth-card d-flex flex-column">
        <div class="d-flex align-center justify-center h-100">
          <v-card variant="flat" class="w-100 mt-0 px-4 pt-12" max-width="500">
            <v-card-text>
              <h4 class="text-h4 mb-2">重置密码</h4>
              <p class="mb-0">请重新输入您的邮箱地址，然后设置新密码。</p>
            </v-card-text>

            <v-card-text class="pb-0 mb-6">
              <v-form>
                <v-row>
                  <v-col cols="12">
                    <v-text-field
                      type="email"
                      autocomplete="email"
                      :autofocus="true"
                      :disabled="updating"
                      label="邮箱"
                      placeholder="请输入邮箱地址"
                      v-model="email"
                      @keyup.enter="passwordInput?.focus()"
                    />
                  </v-col>

                  <v-col cols="12">
                    <v-text-field
                      autocomplete="new-password"
                      ref="passwordInput"
                      type="password"
                      :disabled="updating"
                      label="新密码"
                      placeholder="请输入新密码"
                      v-model="newPassword"
                      @keyup.enter="confirmPasswordInput?.focus()"
                    />
                  </v-col>

                  <v-col cols="12">
                    <v-text-field
                      ref="confirmPasswordInput"
                      type="password"
                      :disabled="updating"
                      label="确认密码"
                      placeholder="请再次输入密码"
                      v-model="confirmPassword"
                      @keyup.enter="resetPassword"
                    />
                  </v-col>

                  <v-col cols="12">
                    <v-btn block :disabled="!email || !newPassword || !confirmPassword || updating" @click="resetPassword">
                      更新密码
                      <v-progress-circular indeterminate size="22" class="ms-2" v-if="updating"></v-progress-circular>
                    </v-btn>
                  </v-col>

                  <v-col cols="12">
                    <router-link class="d-flex align-center justify-center" to="/login"
                                 :class="{ 'disabled': updating }">
                      <v-icon class="icon-with-direction" :icon="mdiChevronLeft"/>
                      <span>返回登录页面</span>
                    </router-link>
                  </v-col>
                </v-row>
              </v-form>
            </v-card-text>
          </v-card>
        </div>
        <v-spacer/>
        <div class="d-flex align-center justify-center">
          <v-card variant="flat" class="w-100 px-4 pb-4" max-width="500">
            <v-card-text class="pt-0">
              <v-row>
                <v-col cols="12" class="text-center">
                  <language-select-button :disabled="updating" />
                </v-col>

                <v-col cols="12" class="d-flex align-center pt-0">
                  <v-divider />
                </v-col>

                <v-col cols="12" class="text-center text-sm">
                  <span>Powered by </span>
                  <a href="https://github.com/zzhao/OmniEvent" target="_blank">OmniEvent</a>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </div>
      </v-col>
    </v-row>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" @update:show="onSnackbarShowStateChanged" />
  </div>
</template>

<script setup lang="ts">
import { VTextField } from 'vuetify/components/VTextField';
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';
import { useTheme } from 'vuetify';

import { useUserStore } from '@/stores/user';

import { ThemeType } from '@/core/theme';
import { APPLICATION_LOGO_PATH } from '@/consts/asset';

import {
    mdiChevronLeft
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const props = defineProps<{
    token: string;
}>();

const router = useRouter();
const theme = useTheme();

const userStore = useUserStore();

const passwordInput = useTemplateRef<VTextField>('passwordInput');
const confirmPasswordInput = useTemplateRef<VTextField>('confirmPasswordInput');
const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const email = ref<string>('');
const newPassword = ref<string>('');
const confirmPassword = ref<string>('');
const updating = ref<boolean>(false);
const passwordChanged = ref<boolean>(false);

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

const inputProblemMessage = computed<string | null>(() => {
    if (!email.value) {
        return '邮箱地址不能为空';
    } else if (!newPassword.value && !confirmPassword.value) {
        return '密码未被修改';
    } else if (!newPassword.value && confirmPassword.value) {
        return '新密码不能为空';
    } else if (newPassword.value && !confirmPassword.value) {
        return '确认密码不能为空';
    } else if (newPassword.value && confirmPassword.value && newPassword.value !== confirmPassword.value) {
        return '两次输入的密码不一致';
    } else {
        return null;
    }
});

function onSnackbarShowStateChanged(newValue: boolean): void {
    if (!newValue && passwordChanged.value) {
        router.replace('/login');
    }
}

function resetPassword(): void  {
    passwordChanged.value = false;

    const problemMessage = inputProblemMessage.value;

    if (problemMessage) {
        snackbar.value?.showMessage(problemMessage);
        return;
    }

    updating.value = true;

    userStore.resetPassword({
        token: props.token,
        email: email.value,
        password: newPassword.value
    }).then(() => {
        updating.value = false;
        passwordChanged.value = true;
        snackbar.value?.showMessage('密码已更新');
    }).catch(error => {
        updating.value = false;
        passwordChanged.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}
</script>

<style scoped>
.auth-logo {
  padding: 24px 40px;
}

.login-page-logo {
  width: 40px;
  height: 40px;
}

.auth-wrapper {
  height: calc(100vh - 88px);
}

.auth-image-background {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100%;
}

.auth-card {
  background: rgb(var(--v-theme-surface));
}

.auth-img-footer {
  position: absolute;
  bottom: 0;
  width: 100%;
}

.img-with-direction {
  width: 100%;
}

.gap-x-3 {
  gap: 12px;
}

.icon-with-direction {
  direction: rtl;
}
</style>