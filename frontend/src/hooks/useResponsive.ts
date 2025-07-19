import { useState, useEffect, useCallback } from 'react';

export interface BreakpointConfig {
  mobile: number;
  tablet: number;
  desktop: number;
  large: number;
}

export interface ViewportInfo {
  width: number;
  height: number;
  isMobile: boolean;
  isTablet: boolean;
  isDesktop: boolean;
  isLarge: boolean;
  orientation: 'portrait' | 'landscape';
  touchDevice: boolean;
  breakpoint: 'mobile' | 'tablet' | 'desktop' | 'large';
}

const DEFAULT_BREAKPOINTS: BreakpointConfig = {
  mobile: 768,
  tablet: 1024,
  desktop: 1200,
  large: 1440
};

export const useResponsive = (customBreakpoints?: Partial<BreakpointConfig>) => {
  const breakpoints = { ...DEFAULT_BREAKPOINTS, ...customBreakpoints };
  
  const [viewportInfo, setViewportInfo] = useState<ViewportInfo>(() => {
    if (typeof window === 'undefined') {
      return {
        width: 1024,
        height: 768,
        isMobile: false,
        isTablet: true,
        isDesktop: false,
        isLarge: false,
        orientation: 'landscape',
        touchDevice: false,
        breakpoint: 'tablet'
      };
    }

    const width = window.innerWidth;
    const height = window.innerHeight;
    const isMobile = width < breakpoints.mobile;
    const isTablet = width >= breakpoints.mobile && width < breakpoints.tablet;
    const isDesktop = width >= breakpoints.tablet && width < breakpoints.large;
    const isLarge = width >= breakpoints.large;
    const orientation = width > height ? 'landscape' : 'portrait';
    const touchDevice = 'ontouchstart' in window || navigator.maxTouchPoints > 0;

    let breakpoint: ViewportInfo['breakpoint'] = 'desktop';
    if (isMobile) breakpoint = 'mobile';
    else if (isTablet) breakpoint = 'tablet';
    else if (isLarge) breakpoint = 'large';

    return {
      width,
      height,
      isMobile,
      isTablet,
      isDesktop,
      isLarge,
      orientation,
      touchDevice,
      breakpoint
    };
  });

  const updateViewportInfo = useCallback(() => {
    const width = window.innerWidth;
    const height = window.innerHeight;
    const isMobile = width < breakpoints.mobile;
    const isTablet = width >= breakpoints.mobile && width < breakpoints.tablet;
    const isDesktop = width >= breakpoints.tablet && width < breakpoints.large;
    const isLarge = width >= breakpoints.large;
    const orientation = width > height ? 'landscape' : 'portrait';
    const touchDevice = 'ontouchstart' in window || navigator.maxTouchPoints > 0;

    let breakpoint: ViewportInfo['breakpoint'] = 'desktop';
    if (isMobile) breakpoint = 'mobile';
    else if (isTablet) breakpoint = 'tablet';
    else if (isLarge) breakpoint = 'large';

    setViewportInfo({
      width,
      height,
      isMobile,
      isTablet,
      isDesktop,
      isLarge,
      orientation,
      touchDevice,
      breakpoint
    });
  }, [breakpoints]);

  useEffect(() => {
    let timeoutId: NodeJS.Timeout;

    const handleResize = () => {
      clearTimeout(timeoutId);
      timeoutId = setTimeout(updateViewportInfo, 100);
    };

    window.addEventListener('resize', handleResize);
    window.addEventListener('orientationchange', updateViewportInfo);

    return () => {
      window.removeEventListener('resize', handleResize);
      window.removeEventListener('orientationchange', updateViewportInfo);
      clearTimeout(timeoutId);
    };
  }, [updateViewportInfo]);

  const isBreakpoint = useCallback((breakpoint: keyof BreakpointConfig) => {
    return viewportInfo.breakpoint === breakpoint;
  }, [viewportInfo.breakpoint]);

  const isBreakpointUp = useCallback((breakpoint: keyof BreakpointConfig) => {
    const currentWidth = viewportInfo.width;
    return currentWidth >= breakpoints[breakpoint];
  }, [viewportInfo.width, breakpoints]);

  const isBreakpointDown = useCallback((breakpoint: keyof BreakpointConfig) => {
    const currentWidth = viewportInfo.width;
    return currentWidth < breakpoints[breakpoint];
  }, [viewportInfo.width, breakpoints]);

  const useBreakpointValue = useCallback(<T>(values: {
    mobile?: T;
    tablet?: T;
    desktop?: T;
    large?: T;
    default: T;
  }): T => {
    const { isMobile, isTablet, isDesktop, isLarge } = viewportInfo;
    
    if (isMobile && values.mobile !== undefined) return values.mobile;
    if (isTablet && values.tablet !== undefined) return values.tablet;
    if (isDesktop && values.desktop !== undefined) return values.desktop;
    if (isLarge && values.large !== undefined) return values.large;
    
    return values.default;
  }, [viewportInfo]);

  const getResponsiveClasses = useCallback((classes: {
    mobile?: string;
    tablet?: string;
    desktop?: string;
    large?: string;
    default?: string;
  }): string => {
    const classNames: string[] = [];
    
    if (classes.default) classNames.push(classes.default);
    if (classes.mobile) classNames.push(classes.mobile);
    if (classes.tablet) classNames.push(classes.tablet);
    if (classes.desktop) classNames.push(classes.desktop);
    if (classes.large) classNames.push(classes.large);
    
    return classNames.join(' ');
  }, []);

  const matchMediaQuery = useCallback((query: string): boolean => {
    if (typeof window === 'undefined') return false;
    return window.matchMedia(query).matches;
  }, []);

  // Predefined media queries
  const mediaQueries = {
    mobile: `(max-width: ${breakpoints.mobile - 1}px)`,
    tablet: `(min-width: ${breakpoints.mobile}px) and (max-width: ${breakpoints.tablet - 1}px)`,
    desktop: `(min-width: ${breakpoints.tablet}px) and (max-width: ${breakpoints.large - 1}px)`,
    large: `(min-width: ${breakpoints.large}px)`,
    mobileUp: `(min-width: ${breakpoints.mobile}px)`,
    tabletUp: `(min-width: ${breakpoints.tablet}px)`,
    desktopUp: `(min-width: ${breakpoints.large}px)`,
    mobileDown: `(max-width: ${breakpoints.mobile - 1}px)`,
    tabletDown: `(max-width: ${breakpoints.tablet - 1}px)`,
    desktopDown: `(max-width: ${breakpoints.large - 1}px)`,
    portrait: '(orientation: portrait)',
    landscape: '(orientation: landscape)',
    retina: '(-webkit-min-device-pixel-ratio: 2), (min-resolution: 192dpi)',
    touch: '(hover: none) and (pointer: coarse)',
    reducedMotion: '(prefers-reduced-motion: reduce)',
    darkMode: '(prefers-color-scheme: dark)',
    lightMode: '(prefers-color-scheme: light)'
  };

  const useMediaQuery = useCallback((query: keyof typeof mediaQueries | string) => {
    const [matches, setMatches] = useState(() => {
      if (typeof window === 'undefined') return false;
      const mediaQuery = typeof query === 'string' && query.startsWith('(') 
        ? query 
        : mediaQueries[query as keyof typeof mediaQueries];
      return window.matchMedia(mediaQuery).matches;
    });

    useEffect(() => {
      const mediaQuery = typeof query === 'string' && query.startsWith('(') 
        ? query 
        : mediaQueries[query as keyof typeof mediaQueries];
      
      const mediaQueryList = window.matchMedia(mediaQuery);
      const handler = (e: MediaQueryListEvent) => setMatches(e.matches);
      
      mediaQueryList.addEventListener('change', handler);
      setMatches(mediaQueryList.matches);
      
      return () => mediaQueryList.removeEventListener('change', handler);
    }, [query]);

    return matches;
  }, []);

  // Layout helpers
  const shouldCollapseSidebar = useCallback(() => {
    return viewportInfo.isMobile;
  }, [viewportInfo.isMobile]);

  const shouldShowMobileMenu = useCallback(() => {
    return viewportInfo.isMobile;
  }, [viewportInfo.isMobile]);

  const getOptimalColumns = useCallback((
    itemWidth: number,
    containerWidth?: number,
    gap: number = 16
  ): number => {
    const availableWidth = containerWidth || viewportInfo.width;
    const totalItemWidth = itemWidth + gap;
    const columns = Math.floor((availableWidth + gap) / totalItemWidth);
    return Math.max(1, columns);
  }, [viewportInfo.width]);

  const getTouchTargetSize = useCallback(() => {
    return viewportInfo.touchDevice ? 44 : 32;
  }, [viewportInfo.touchDevice]);

  return {
    // Viewport information
    ...viewportInfo,
    
    // Breakpoint utilities
    isBreakpoint,
    isBreakpointUp,
    isBreakpointDown,
    
    // Responsive values
    useBreakpointValue,
    getResponsiveClasses,
    
    // Media queries
    matchMediaQuery,
    useMediaQuery,
    mediaQueries,
    
    // Layout helpers
    shouldCollapseSidebar,
    shouldShowMobileMenu,
    getOptimalColumns,
    getTouchTargetSize,
    
    // Configuration
    breakpoints
  };
};

export default useResponsive;