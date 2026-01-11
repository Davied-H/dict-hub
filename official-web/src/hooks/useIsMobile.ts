import { useState, useEffect } from 'react';

const MOBILE_BREAKPOINT = 768;

/**
 * 检测当前设备是否为移动端
 * 用于在移动端禁用或简化复杂动画以提升性能
 */
export function useIsMobile(): boolean {
  const [isMobile, setIsMobile] = useState(() => {
    // SSR 安全检查
    if (typeof window === 'undefined') return false;
    return window.innerWidth < MOBILE_BREAKPOINT;
  });

  useEffect(() => {
    // 使用 matchMedia 以获得更好的性能
    const mediaQuery = window.matchMedia(`(max-width: ${MOBILE_BREAKPOINT - 1}px)`);

    const handleChange = (e: MediaQueryListEvent | MediaQueryList) => {
      setIsMobile(e.matches);
    };

    // 初始化
    handleChange(mediaQuery);

    // 监听变化
    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, []);

  return isMobile;
}

/**
 * 检测是否应该减少动画
 * 考虑用户的系统偏好设置（prefers-reduced-motion）
 */
export function useReducedMotion(): boolean {
  const [prefersReduced, setPrefersReduced] = useState(() => {
    if (typeof window === 'undefined') return false;
    return window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  });

  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)');

    const handleChange = (e: MediaQueryListEvent) => {
      setPrefersReduced(e.matches);
    };

    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, []);

  return prefersReduced;
}

/**
 * 综合判断是否应该使用简化动画
 * 移动端或用户偏好减少动画时返回 true
 */
export function useShouldReduceMotion(): boolean {
  const isMobile = useIsMobile();
  const prefersReduced = useReducedMotion();
  return isMobile || prefersReduced;
}
