import { AccountSettingsOperationStatus } from '../proto/byte/v/forge/waapp/v1/account_settings';
import type { GetTwoFactorAuthStatusResponse, TwoFactorAuthStatus } from '../proto/byte/v/forge/waapp/v1/account_settings';
import { i18n } from '@/i18n/i18n';

type BadgeVariant = 'default' | 'secondary' | 'destructive' | 'outline';

export type TwoFactorStatusView = { isFetching: boolean; isError: boolean; data?: { status?: TwoFactorAuthStatus } };

export function initialTwoFactorStatus(status?: TwoFactorAuthStatus): GetTwoFactorAuthStatusResponse {
  return status ? { status, error: undefined } : { status: undefined, error: undefined };
}

export function shouldCollectEmailOtpAfterSet(status?: AccountSettingsOperationStatus) {
  return status !== AccountSettingsOperationStatus.ACCOUNT_SETTINGS_OPERATION_STATUS_VERIFIED
    && status !== AccountSettingsOperationStatus.ACCOUNT_SETTINGS_OPERATION_STATUS_REJECTED;
}

export function shouldShowEmailOtp(status?: AccountSettingsOperationStatus) {
  return status === AccountSettingsOperationStatus.ACCOUNT_SETTINGS_OPERATION_STATUS_NEEDS_VERIFICATION
    || status === AccountSettingsOperationStatus.ACCOUNT_SETTINGS_OPERATION_STATUS_WAITING
    || status === AccountSettingsOperationStatus.ACCOUNT_SETTINGS_OPERATION_STATUS_CODE_MISMATCH;
}

export function twoFactorStatusLabel(query: TwoFactorStatusView) {
  if (query.isFetching) return i18n.t('account.security.status.syncing', '同步中');
  if (query.isError) return i18n.t('account.security.status.sync_failed', '同步失败');
  if (!query.data?.status) return i18n.t('account.security.status.not_synced', '未同步');
  return query.data.status.configured ? i18n.t('account.security.status.configured', '已配置') : i18n.t('account.security.status.not_configured', '未配置');
}

export function emailStatusLabel(query: TwoFactorStatusView) {
  if (query.isFetching) return i18n.t('account.security.status.syncing', '同步中');
  if (query.isError) return i18n.t('account.security.status.sync_failed', '同步失败');
  if (!query.data?.status) return i18n.t('account.security.status.not_synced', '未同步');
  if (query.data.status.email_verified) return i18n.t('account.security.status.verified', '已验证');
  if (query.data.status.email_address) return i18n.t('account.security.status.pending_verification', '待验证');
  return query.data.status.email_configured ? i18n.t('account.security.status.configured', '已配置') : i18n.t('account.security.status.not_configured', '未配置');
}

export function twoFactorBadgeVariant(query: TwoFactorStatusView): BadgeVariant {
  if (query.isError) return 'destructive';
  return query.data?.status?.configured ? 'default' : 'outline';
}

export function emailBadgeVariant(query: TwoFactorStatusView): BadgeVariant {
  if (query.isError) return 'destructive';
  if (query.data?.status?.email_verified) return 'default';
  return query.data?.status?.email_address || query.data?.status?.email_configured ? 'secondary' : 'outline';
}

export function twoFactorConfigured(query: TwoFactorStatusView) {
  return Boolean(query.data?.status?.configured);
}

export function twoFactorEmailConfigured(query: TwoFactorStatusView) {
  return Boolean(query.data?.status?.email_configured || query.data?.status?.email_address);
}

export function statusLabel(status?: AccountSettingsOperationStatus) {
  switch (status) {
    case AccountSettingsOperationStatus.ACCOUNT_SETTINGS_OPERATION_STATUS_NEEDS_VERIFICATION: return i18n.t('account.security.status.needs_email_verification', '待邮箱验证');
    case AccountSettingsOperationStatus.ACCOUNT_SETTINGS_OPERATION_STATUS_WAITING: return i18n.t('account.security.status.waiting_otp', '等待 OTP');
    case AccountSettingsOperationStatus.ACCOUNT_SETTINGS_OPERATION_STATUS_VERIFIED: return i18n.t('account.security.status.verified', '已验证');
    case AccountSettingsOperationStatus.ACCOUNT_SETTINGS_OPERATION_STATUS_CODE_MISMATCH: return i18n.t('account.security.status.code_mismatch', '验证码不匹配');
    case AccountSettingsOperationStatus.ACCOUNT_SETTINGS_OPERATION_STATUS_REJECTED: return i18n.t('account.security.status.rejected', '已拒绝');
    case AccountSettingsOperationStatus.ACCOUNT_SETTINGS_OPERATION_STATUS_ACCEPTED: return i18n.t('account.security.status.accepted', '已受理');
    default: return i18n.t('account.security.status.not_started', '未执行');
  }
}
