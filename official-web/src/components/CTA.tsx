import { motion, useScroll, useTransform } from 'framer-motion';
import { useRef } from 'react';
import { useShouldReduceMotion } from '../hooks/useIsMobile';

export default function CTA() {
  const containerRef = useRef<HTMLDivElement>(null);
  const reduceMotion = useShouldReduceMotion();

  const { scrollYProgress } = useScroll({
    target: containerRef,
    offset: ['start end', 'end start'],
  });

  const y = useTransform(scrollYProgress, [0, 1], [50, -50]);
  const scale = useTransform(scrollYProgress, [0, 0.5], [0.95, 1]);

  // 移动端简化版本
  if (reduceMotion) {
    return (
      <section
        ref={containerRef}
        className="relative py-32 overflow-hidden"
      >
        {/* Static gradient background */}
        <div className="absolute inset-0 bg-gradient-primary" />

        <div className="relative max-w-3xl mx-auto px-6 text-center">
          <h2 className="text-4xl lg:text-6xl font-extrabold text-white mb-6">
            开始使用 Dict-Hub
          </h2>

          <p className="text-xl text-white/80 mb-12">
            立即体验现代化的词典查询服务，提升你的语言学习效率
          </p>

          <div className="flex flex-wrap justify-center gap-4">
            <a
              href="https://github.com/Davied-H/dict-hub"
              target="_blank"
              rel="noopener noreferrer"
              className="group flex items-center gap-3 px-8 py-4 bg-white text-primary rounded-2xl font-bold text-lg shadow-xl"
            >
              <svg className="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
              </svg>
              在 GitHub 上查看
              <span className="opacity-70">→</span>
            </a>

            <a
              href="#start"
              className="flex items-center gap-3 px-8 py-4 bg-white/10 text-white rounded-2xl font-bold text-lg border border-white/20 hover:bg-white/20 transition-colors"
            >
              快速开始
            </a>
          </div>
        </div>
      </section>
    );
  }

  // 桌面端完整动画版本
  return (
    <section
      ref={containerRef}
      className="relative py-32 overflow-hidden"
    >
      {/* Animated gradient background */}
      <motion.div
        className="absolute inset-0 bg-gradient-primary"
        style={{ y }}
      >
        {/* Animated shapes */}
        <motion.div
          className="absolute top-0 left-0 w-[500px] h-[500px] bg-white/10 rounded-full blur-3xl"
          animate={{
            x: [-200, 100, -200],
            y: [-100, 200, -100],
          }}
          transition={{ duration: 15, repeat: Infinity, ease: 'easeInOut' }}
        />
        <motion.div
          className="absolute bottom-0 right-0 w-[600px] h-[600px] bg-white/10 rounded-full blur-3xl"
          animate={{
            x: [200, -100, 200],
            y: [100, -200, 100],
          }}
          transition={{ duration: 18, repeat: Infinity, ease: 'easeInOut' }}
        />
      </motion.div>

      <motion.div
        className="relative max-w-3xl mx-auto px-6 text-center"
        style={{ scale }}
      >
        <motion.h2
          className="text-4xl lg:text-6xl font-extrabold text-white mb-6"
          initial={{ opacity: 0, y: 40 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.8, ease: [0.16, 1, 0.3, 1] }}
        >
          开始使用 Dict-Hub
        </motion.h2>

        <motion.p
          className="text-xl text-white/80 mb-12"
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.1, duration: 0.8 }}
        >
          立即体验现代化的词典查询服务，提升你的语言学习效率
        </motion.p>

        <motion.div
          className="flex flex-wrap justify-center gap-4"
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.2, duration: 0.8 }}
        >
          <motion.a
            href="https://github.com/Davied-H/dict-hub"
            target="_blank"
            rel="noopener noreferrer"
            className="group flex items-center gap-3 px-8 py-4 bg-white text-primary rounded-2xl font-bold text-lg shadow-xl"
            whileHover={{ scale: 1.05, y: -4 }}
            whileTap={{ scale: 0.95 }}
          >
            <svg className="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
            </svg>
            在 GitHub 上查看
            <motion.span
              className="opacity-0 group-hover:opacity-100 transition-opacity"
              animate={{ x: [0, 4, 0] }}
              transition={{ duration: 1.5, repeat: Infinity }}
            >
              →
            </motion.span>
          </motion.a>

          <motion.a
            href="#start"
            className="flex items-center gap-3 px-8 py-4 bg-white/10 text-white rounded-2xl font-bold text-lg border border-white/20 hover:bg-white/20 transition-colors"
            whileHover={{ scale: 1.05, y: -4 }}
            whileTap={{ scale: 0.95 }}
          >
            快速开始
          </motion.a>
        </motion.div>
      </motion.div>
    </section>
  );
}
