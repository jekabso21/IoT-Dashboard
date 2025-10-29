import { CheckCircle, XCircle, AlertTriangle, Info, X } from 'lucide-react';

interface ToastProps {
  id: string;
  type: 'success' | 'error' | 'warning' | 'info';
  message: string;
  onClose: () => void;
}

export default function Toast({ type, message, onClose }: ToastProps) {
  const getVariantStyles = () => {
    switch (type) {
      case 'success':
        return {
          container: 'bg-elevated border-l-4 border-green-500',
          icon: <CheckCircle className="w-5 h-5 text-green-500" />,
        };
      case 'error':
        return {
          container: 'bg-elevated border-l-4 border-red-500',
          icon: <XCircle className="w-5 h-5 text-red-500" />,
        };
      case 'warning':
        return {
          container: 'bg-elevated border-l-4 border-orange-500',
          icon: <AlertTriangle className="w-5 h-5 text-orange-500" />,
        };
      case 'info':
        return {
          container: 'bg-elevated border-l-4 border-accent',
          icon: <Info className="w-5 h-5 text-accent" />,
        };
      default:
        return {
          container: 'bg-elevated border-l-4 border-accent',
          icon: <Info className="w-5 h-5 text-accent" />,
        };
    }
  };

  const { container, icon } = getVariantStyles();

  return (
    <div 
      className={`${container} flex items-center gap-3 p-4 rounded-lg shadow-lg toast-container transition-all duration-300`}
      role="status"
      aria-live="polite"
    >
      {icon}
      <span className="flex-1 text-primary">{message}</span>
      <button
        onClick={onClose}
        className="ml-auto cursor-pointer text-secondary hover:text-primary transition-colors"
        aria-label="Close notification"
        type="button"
      >
        <X className="w-4 h-4" />
      </button>
    </div>
  );
}
