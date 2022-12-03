import { Store } from 'vuex';
import { NavigationGuard } from 'vue-router';
import { RootState, State } from '@/store';

export function authorize(store: Store<RootState>): NavigationGuard {
  return (to, from, next) => {
    if (to.meta && to.meta.requiresAuth && !(store.state as State).user.current) {
      window.location.href = `/login?redirect=${encodeURIComponent(to.fullPath)}`;
    } else {
      next();
    }
  };
}
