import { useEffect, useCallback, useRef, useState } from 'react';

export interface AnimationConfig {
  duration?: number;
  delay?: number;
  easing?: string;
  fillMode?: 'none' | 'forwards' | 'backwards' | 'both';
  direction?: 'normal' | 'reverse' | 'alternate' | 'alternate-reverse';
  iterationCount?: number | 'infinite';
}

export interface TransitionConfig {
  property?: string;
  duration?: number;
  easing?: string;
  delay?: number;
}

export const useAnimations = () => {
  const [reducedMotion, setReducedMotion] = useState(false);
  const animatingElementsRef = useRef<Set<HTMLElement>>(new Set());

  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)');
    setReducedMotion(mediaQuery.matches);

    const handleChange = (e: MediaQueryListEvent) => {
      setReducedMotion(e.matches);
    };

    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, []);

  const animate = useCallback((
    element: HTMLElement,
    keyframes: Keyframe[] | PropertyIndexedKeyframes,
    options: KeyframeAnimationOptions & AnimationConfig = {}
  ): Animation | null => {
    if (reducedMotion) return null;

    const {
      duration = 300,
      delay = 0,
      easing = 'ease-out',
      fillMode = 'both',
      direction = 'normal',
      iterationCount = 1,
      ...keyframeOptions
    } = options;

    try {
      animatingElementsRef.current.add(element);
      
      const animation = element.animate(keyframes, {
        duration,
        delay,
        easing,
        fill: fillMode,
        direction,
        iterations: iterationCount,
        ...keyframeOptions
      });

      animation.addEventListener('finish', () => {
        animatingElementsRef.current.delete(element);
      });

      animation.addEventListener('cancel', () => {
        animatingElementsRef.current.delete(element);
      });

      return animation;
    } catch (error) {
      console.warn('Animation failed:', error);
      return null;
    }
  }, [reducedMotion]);

  const fadeIn = useCallback((
    element: HTMLElement,
    options: AnimationConfig = {}
  ) => {
    return animate(element, [
      { opacity: 0 },
      { opacity: 1 }
    ], { duration: 300, ...options });
  }, [animate]);

  const fadeOut = useCallback((
    element: HTMLElement,
    options: AnimationConfig = {}
  ) => {
    return animate(element, [
      { opacity: 1 },
      { opacity: 0 }
    ], { duration: 300, ...options });
  }, [animate]);

  const slideIn = useCallback((
    element: HTMLElement,
    direction: 'left' | 'right' | 'up' | 'down' = 'left',
    options: AnimationConfig = {}
  ) => {
    const transforms = {
      left: ['translateX(-100%)', 'translateX(0)'],
      right: ['translateX(100%)', 'translateX(0)'],
      up: ['translateY(-100%)', 'translateY(0)'],
      down: ['translateY(100%)', 'translateY(0)']
    };

    return animate(element, [
      { transform: transforms[direction][0], opacity: 0 },
      { transform: transforms[direction][1], opacity: 1 }
    ], { duration: 400, ...options });
  }, [animate]);

  const slideOut = useCallback((
    element: HTMLElement,
    direction: 'left' | 'right' | 'up' | 'down' = 'right',
    options: AnimationConfig = {}
  ) => {
    const transforms = {
      left: ['translateX(0)', 'translateX(-100%)'],
      right: ['translateX(0)', 'translateX(100%)'],
      up: ['translateY(0)', 'translateY(-100%)'],
      down: ['translateY(0)', 'translateY(100%)']
    };

    return animate(element, [
      { transform: transforms[direction][0], opacity: 1 },
      { transform: transforms[direction][1], opacity: 0 }
    ], { duration: 400, ...options });
  }, [animate]);

  const scaleIn = useCallback((
    element: HTMLElement,
    options: AnimationConfig = {}
  ) => {
    return animate(element, [
      { transform: 'scale(0.8)', opacity: 0 },
      { transform: 'scale(1)', opacity: 1 }
    ], { duration: 250, easing: 'cubic-bezier(0.34, 1.56, 0.64, 1)', ...options });
  }, [animate]);

  const scaleOut = useCallback((
    element: HTMLElement,
    options: AnimationConfig = {}
  ) => {
    return animate(element, [
      { transform: 'scale(1)', opacity: 1 },
      { transform: 'scale(0.8)', opacity: 0 }
    ], { duration: 250, ...options });
  }, [animate]);

  const bounceIn = useCallback((
    element: HTMLElement,
    options: AnimationConfig = {}
  ) => {
    return animate(element, [
      { transform: 'scale(0.3)', opacity: 0 },
      { transform: 'scale(1.05)', opacity: 0.8, offset: 0.5 },
      { transform: 'scale(0.9)', opacity: 0.9, offset: 0.7 },
      { transform: 'scale(1)', opacity: 1 }
    ], { duration: 600, easing: 'cubic-bezier(0.68, -0.55, 0.265, 1.55)', ...options });
  }, [animate]);

  const shake = useCallback((
    element: HTMLElement,
    options: AnimationConfig = {}
  ) => {
    return animate(element, [
      { transform: 'translateX(0)' },
      { transform: 'translateX(-10px)' },
      { transform: 'translateX(10px)' },
      { transform: 'translateX(-10px)' },
      { transform: 'translateX(10px)' },
      { transform: 'translateX(-5px)' },
      { transform: 'translateX(5px)' },
      { transform: 'translateX(0)' }
    ], { duration: 600, ...options });
  }, [animate]);

  const pulse = useCallback((
    element: HTMLElement,
    options: AnimationConfig = {}
  ) => {
    return animate(element, [
      { opacity: 1 },
      { opacity: 0.5 },
      { opacity: 1 }
    ], { duration: 1000, iterationCount: 'infinite', ...options });
  }, [animate]);

  const spin = useCallback((
    element: HTMLElement,
    options: AnimationConfig = {}
  ) => {
    return animate(element, [
      { transform: 'rotate(0deg)' },
      { transform: 'rotate(360deg)' }
    ], { duration: 1000, iterationCount: 'infinite', easing: 'linear', ...options });
  }, [animate]);

  const setTransition = useCallback((
    element: HTMLElement,
    config: TransitionConfig
  ) => {
    if (reducedMotion) return;

    const {
      property = 'all',
      duration = 300,
      easing = 'ease-out',
      delay = 0
    } = config;

    element.style.transition = `${property} ${duration}ms ${easing} ${delay}ms`;
  }, [reducedMotion]);

  const clearTransition = useCallback((element: HTMLElement) => {
    element.style.transition = '';
  }, []);

  const staggerChildren = useCallback((
    container: HTMLElement,
    animationFn: (element: HTMLElement, index: number) => Animation | null,
    staggerDelay: number = 100
  ) => {
    const children = Array.from(container.children) as HTMLElement[];
    const animations: (Animation | null)[] = [];

    children.forEach((child, index) => {
      setTimeout(() => {
        const animation = animationFn(child, index);
        animations.push(animation);
      }, index * staggerDelay);
    });

    return animations;
  }, []);

  const animateProgressBar = useCallback((
    element: HTMLElement,
    progress: number,
    options: AnimationConfig = {}
  ) => {
    return animate(element, [
      { width: '0%' },
      { width: `${Math.max(0, Math.min(100, progress))}%` }
    ], { duration: 1000, ...options });
  }, [animate]);

  const morphPath = useCallback((
    svgPath: SVGPathElement,
    newPath: string,
    options: AnimationConfig = {}
  ) => {
    const currentPath = svgPath.getAttribute('d') || '';
    
    return animate(svgPath, [
      { d: currentPath },
      { d: newPath }
    ], { duration: 500, ...options });
  }, [animate]);

  const cancelAllAnimations = useCallback(() => {
    animatingElementsRef.current.forEach(element => {
      const animations = element.getAnimations();
      animations.forEach(animation => animation.cancel());
    });
    animatingElementsRef.current.clear();
  }, []);

  const isAnimating = useCallback((element: HTMLElement): boolean => {
    return animatingElementsRef.current.has(element);
  }, []);

  return {
    // Basic animations
    animate,
    fadeIn,
    fadeOut,
    slideIn,
    slideOut,
    scaleIn,
    scaleOut,
    bounceIn,
    shake,
    pulse,
    spin,
    
    // Transitions
    setTransition,
    clearTransition,
    
    // Advanced animations
    staggerChildren,
    animateProgressBar,
    morphPath,
    
    // Utilities
    cancelAllAnimations,
    isAnimating,
    
    // State
    reducedMotion,
    animatingElements: animatingElementsRef.current
  };
};

export default useAnimations;