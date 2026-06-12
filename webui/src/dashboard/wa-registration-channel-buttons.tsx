import { CheckCircle2, Clock3, CircleDashed, PhoneMissed, XCircle } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { i18n } from '@/i18n/i18n';
import { countdownLabel } from './wa-result-labels';
import {
  registrationMethodAvailable,
  registrationMethodCooldownSeconds,
  registrationMethodStatus,
  visibleRegistrationChannelMethods,
  type RegistrationChannelMethodOption,
  type SelectableRegistrationMethodOption,
} from './wa-registration-methods';
import type { WaProbeStatus } from './wa-result-model';

type Props = {
  status: WaProbeStatus | null;
  elapsedSeconds: number;
  disabled?: boolean;
  onStart: (method: SelectableRegistrationMethodOption) => void;
};

export function WaRegistrationChannelButtons({ status, elapsedSeconds, disabled, onStart }: Props) {
  const methods = status ? visibleRegistrationChannelMethods.filter((method) => registrationMethodStatus(status, method.value)) : visibleRegistrationChannelMethods;
  return (
    <div className="grid gap-2 sm:grid-cols-2">
      {methods.map((method) => {
        const state = channelState(method, status, elapsedSeconds);
        return (
          <Button
            key={method.value}
            type="button"
            variant={state.ready ? 'default' : state.cooldown > 0 ? 'secondary' : 'outline'}
            className="h-10 justify-between gap-3 px-3"
            disabled={disabled || !state.ready}
            aria-label={`${method.label} ${state.label}`}
            title={state.title}
            onClick={() => method.directRequest && onStart(method)}
          >
            <span className="truncate text-sm font-medium">{method.label}</span>
            <Badge variant={state.badge} className="shrink-0">
              <state.Icon />
              {state.label}
            </Badge>
          </Button>
        );
      })}
    </div>
  );
}

function channelState(method: RegistrationChannelMethodOption, status: WaProbeStatus | null, elapsedSeconds: number) {
  if (!status) return { ready: false, cooldown: 0, label: i18n.t('registration.channel.probe_first', '先检测'), badge: 'outline' as const, Icon: CircleDashed, title: i18n.t('registration.channel.probe_first', '先检测') };
  const cooldown = registrationMethodCooldownSeconds(status, method.value, elapsedSeconds);
  if (cooldown > 0) {
    return { ready: false, cooldown, label: countdownLabel(cooldown), badge: 'secondary' as const, Icon: Clock3, title: i18n.t('registration.channel.cooling_down', '冷却中') };
  }
  if (!method.directRequest) {
    return { ready: false, cooldown: 0, label: i18n.t('registration.channel.unsupported', '不支持'), badge: 'outline' as const, Icon: PhoneMissed, title: i18n.t('registration.channel.unsupported', '不支持') };
  }
  if (registrationMethodAvailable(status, method.value, elapsedSeconds)) {
    return { ready: true, cooldown: 0, label: i18n.t('registration.channel.available', '可用'), badge: 'default' as const, Icon: CheckCircle2, title: i18n.t('registration.channel.available', '可用') };
  }
  return { ready: false, cooldown: 0, label: i18n.t('registration.channel.unavailable', '不可用'), badge: 'outline' as const, Icon: XCircle, title: `${method.label} ${i18n.t('registration.channel.currently_unavailable', '当前不可用')}` };
}
