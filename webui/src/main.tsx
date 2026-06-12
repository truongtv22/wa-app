import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { createRoot } from 'react-dom/client';
import { RouterProvider } from 'react-router/dom';
import { TooltipProvider } from '@/components/ui/tooltip';
import { initI18n } from '@/i18n/i18n';
import { waRouter } from './dashboard/wa-routes';
import './styles.css';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { refetchOnWindowFocus: false, retry: 1, staleTime: 5000 },
  },
});

initI18n();

createRoot(document.getElementById('root')!).render(
  <QueryClientProvider client={queryClient}>
    <TooltipProvider>
      <RouterProvider router={waRouter} />
    </TooltipProvider>
  </QueryClientProvider>,
);
