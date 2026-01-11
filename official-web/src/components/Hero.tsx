import { motion, useScroll, useTransform } from 'framer-motion';
import { useRef } from 'react';

export default function Hero() {
  const containerRef = useRef<HTMLDivElement>(null);
  const { scrollYProgress } = useScroll({
    target: containerRef,
    offset: ['start start', 'end start'],
  });

  const y = useTransform(scrollYProgress, [0, 1], ['0%', '50%']);
  const opacity = useTransform(scrollYProgress, [0, 0.5], [1, 0]);
  const scale = useTransform(scrollYProgress, [0, 0.5], [1, 0.9]);
  const textY = useTransform(scrollYProgress, [0, 1], ['0%', '100%']);

  const cardRotateX = useTransform(scrollYProgress, [0, 0.5], [0, 10]);
  const cardScale = useTransform(scrollYProgress, [0, 0.5], [1, 0.95]);

  return (
    <section
      ref={containerRef}
      className="relative min-h-[120vh] flex items-center overflow-hidden bg-gradient-to-b from-slate-50 to-white"
    >
      {/* Background animated elements */}
      <motion.div
        className="absolute inset-0 overflow-hidden"
        style={{ opacity }}
      >
        <motion.div
          className="absolute top-1/4 -left-32 w-96 h-96 bg-primary/10 rounded-full blur-3xl"
          animate={{
            x: [0, 50, 0],
            y: [0, 30, 0],
          }}
          transition={{
            duration: 8,
            repeat: Infinity,
            ease: 'easeInOut',
          }}
        />
        <motion.div
          className="absolute top-1/3 -right-32 w-96 h-96 bg-accent/10 rounded-full blur-3xl"
          animate={{
            x: [0, -50, 0],
            y: [0, -30, 0],
          }}
          transition={{
            duration: 10,
            repeat: Infinity,
            ease: 'easeInOut',
          }}
        />
      </motion.div>

      <div className="relative max-w-7xl mx-auto px-6 pt-32 pb-20">
        <motion.div
          className="grid lg:grid-cols-2 gap-16 items-center"
          style={{ y, opacity, scale }}
        >
          {/* Left content */}
          <motion.div style={{ y: textY }}>
            <motion.div
              initial={{ opacity: 0, y: 40 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.8, ease: [0.16, 1, 0.3, 1] }}
            >
              <motion.span
                className="inline-block px-4 py-2 bg-primary/10 text-primary rounded-full text-sm font-semibold mb-6"
                initial={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ delay: 0.2 }}
              >
                开源 MDX 词典查询系统
              </motion.span>

              <h1 className="text-5xl lg:text-7xl font-extrabold leading-tight mb-8">
                <span className="text-gradient">现代化</span>
                <br />
                <span className="text-slate-900">词典查询平台</span>
              </h1>

              <p className="text-xl text-slate-600 leading-relaxed mb-10 max-w-lg">
                Dict-Hub 是一个高性能的 MDX 词典查询系统，支持多词典管理、快速检索、历史记录追踪与词频统计。
              </p>

              <div className="flex flex-wrap gap-4">
                <motion.a
                  href="#start"
                  className="group flex items-center gap-3 px-8 py-4 bg-gradient-primary text-white rounded-2xl font-semibold text-lg shadow-xl shadow-primary/30"
                  whileHover={{ scale: 1.05, y: -4 }}
                  whileTap={{ scale: 0.95 }}
                >
                  快速开始
                  <motion.svg
                    className="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    animate={{ x: [0, 4, 0] }}
                    transition={{ duration: 1.5, repeat: Infinity }}
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M17 8l4 4m0 0l-4 4m4-4H3"
                    />
                  </motion.svg>
                </motion.a>

                <motion.a
                  href="https://github.com/Davied-H/dict-hub"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-3 px-8 py-4 bg-slate-100 text-slate-900 rounded-2xl font-semibold text-lg hover:bg-slate-200 transition-colors"
                  whileHover={{ scale: 1.05, y: -4 }}
                  whileTap={{ scale: 0.95 }}
                >
                  查看文档
                </motion.a>
              </div>
            </motion.div>
          </motion.div>

          {/* Right card - search demo */}
          <motion.div
            className="relative lg:order-first"
            style={{
              rotateX: cardRotateX,
              scale: cardScale,
              transformPerspective: 1000,
            }}
            initial={{ opacity: 0, y: 60, rotateY: -10 }}
            animate={{ opacity: 1, y: 0, rotateY: 0 }}
            transition={{ duration: 1, delay: 0.3, ease: [0.16, 1, 0.3, 1] }}
          >
            <div className="relative bg-white rounded-3xl shadow-2xl border border-slate-200/50 p-8 overflow-hidden">
              {/* Floating decorative elements */}
              <div className="absolute -top-10 -right-10 w-40 h-40 bg-gradient-to-br from-primary/20 to-accent/20 rounded-full blur-2xl" />
              <div className="absolute -bottom-10 -left-10 w-32 h-32 bg-accent/10 rounded-full blur-2xl" />

              {/* Search bar */}
              <motion.div
                className="relative flex items-center gap-4 px-6 py-4 bg-slate-50 rounded-2xl mb-6"
                whileHover={{ scale: 1.02 }}
              >
                <svg
                  className="w-5 h-5 text-slate-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <circle cx="11" cy="11" r="8" strokeWidth={2} />
                  <path
                    strokeLinecap="round"
                    strokeWidth={2}
                    d="M21 21l-4.35-4.35"
                  />
                </svg>
                <motion.span
                  className="text-slate-900 font-medium"
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  transition={{ delay: 0.8 }}
                >
                  eloquent
                  <motion.span
                    className="inline-block w-0.5 h-5 bg-primary ml-1"
                    animate={{ opacity: [1, 0] }}
                    transition={{ duration: 0.8, repeat: Infinity }}
                  />
                </motion.span>
              </motion.div>

              {/* Result items */}
              {[
                {
                  word: 'eloquent',
                  phonetic: '/ˈeləkwənt/',
                  definition: 'adj. 雄辩的，有口才的；有说服力的',
                },
                {
                  word: 'eloquence',
                  phonetic: '/ˈeləkwəns/',
                  definition: 'n. 口才，雄辩；雄辩术',
                },
              ].map((item, i) => (
                <motion.div
                  key={item.word}
                  className="relative p-5 rounded-2xl bg-slate-50/50 hover:bg-slate-100 transition-colors cursor-pointer mb-3 last:mb-0"
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: 1 + i * 0.15 }}
                  whileHover={{ x: 8, scale: 1.02 }}
                >
                  <div className="text-xl font-bold text-primary mb-1">
                    {item.word}
                  </div>
                  <div className="text-sm text-slate-500 mb-2">
                    {item.phonetic}
                  </div>
                  <div className="text-slate-600">{item.definition}</div>
                </motion.div>
              ))}
            </div>
          </motion.div>
        </motion.div>
      </div>

      {/* Scroll indicator */}
      <motion.div
        className="absolute bottom-10 left-1/2 -translate-x-1/2"
        style={{ opacity }}
        animate={{ y: [0, 10, 0] }}
        transition={{ duration: 2, repeat: Infinity }}
      >
        <div className="w-6 h-10 border-2 border-slate-300 rounded-full flex justify-center pt-2">
          <motion.div
            className="w-1.5 h-3 bg-slate-400 rounded-full"
            animate={{ y: [0, 12, 0], opacity: [1, 0.3, 1] }}
            transition={{ duration: 2, repeat: Infinity }}
          />
        </div>
      </motion.div>
    </section>
  );
}
