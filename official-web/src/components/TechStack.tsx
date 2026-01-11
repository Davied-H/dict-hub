import { motion, useScroll, useTransform } from 'framer-motion';
import { useRef } from 'react';

const techStack = {
  backend: {
    icon: '⚡',
    title: '后端技术',
    color: 'from-cyan-500 to-blue-600',
    items: [
      'Go 1.24 - 高性能后端语言',
      'Gin - 轻量级 Web 框架',
      'GORM - ORM 数据库操作',
      'SQLite - 轻量级数据库',
      'MDX Parser - 词典格式解析',
    ],
  },
  frontend: {
    icon: '⚛️',
    title: '前端技术',
    color: 'from-blue-400 to-cyan-500',
    items: [
      'React 19 - 前端框架',
      'TypeScript - 类型安全',
      'Vite - 下一代构建工具',
      'Tailwind CSS - 原子化样式',
      'Zustand + React Query - 状态管理',
    ],
  },
};

function TechCard({
  tech,
  index,
}: {
  tech: (typeof techStack)['backend'];
  index: number;
}) {
  const cardRef = useRef<HTMLDivElement>(null);
  const { scrollYProgress } = useScroll({
    target: cardRef,
    offset: ['start end', 'center center'],
  });

  const x = useTransform(
    scrollYProgress,
    [0, 1],
    [index === 0 ? -100 : 100, 0]
  );
  const opacity = useTransform(scrollYProgress, [0, 0.5], [0, 1]);
  const rotateY = useTransform(
    scrollYProgress,
    [0, 1],
    [index === 0 ? -15 : 15, 0]
  );

  return (
    <motion.div
      ref={cardRef}
      style={{ x, opacity, rotateY, transformPerspective: 1200 }}
      className="relative bg-white rounded-3xl border border-slate-200 overflow-hidden"
    >
      {/* Header */}
      <div className="p-8 pb-6">
        <div className="flex items-center gap-4 mb-6">
          <motion.div
            className={`w-14 h-14 rounded-2xl bg-gradient-to-br ${tech.color} flex items-center justify-center text-2xl shadow-lg`}
            whileHover={{ scale: 1.1, rotate: 360 }}
            transition={{ duration: 0.6 }}
          >
            {tech.icon}
          </motion.div>
          <h3 className="text-2xl font-bold text-slate-900">{tech.title}</h3>
        </div>

        <div className="space-y-1">
          {tech.items.map((item, i) => (
            <motion.div
              key={item}
              className="group flex items-center gap-4 py-4 border-b border-slate-100 last:border-0"
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.1 }}
              whileHover={{ x: 8 }}
            >
              <motion.div
                className="w-6 h-6 rounded-full bg-primary/10 flex items-center justify-center"
                whileHover={{ scale: 1.2, backgroundColor: 'rgb(37 99 235)' }}
              >
                <svg
                  className="w-4 h-4 text-primary group-hover:text-white transition-colors"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M5 13l4 4L19 7"
                  />
                </svg>
              </motion.div>
              <span className="text-slate-700 font-medium group-hover:text-primary transition-colors">
                {item}
              </span>
            </motion.div>
          ))}
        </div>
      </div>
    </motion.div>
  );
}

export default function TechStack() {
  const containerRef = useRef<HTMLDivElement>(null);
  const { scrollYProgress } = useScroll({
    target: containerRef,
    offset: ['start end', 'center center'],
  });

  const scale = useTransform(scrollYProgress, [0, 0.5], [0.9, 1]);
  const opacity = useTransform(scrollYProgress, [0, 0.3], [0, 1]);

  return (
    <section
      ref={containerRef}
      id="tech"
      className="relative py-32 bg-white overflow-hidden"
    >
      {/* Animated background */}
      <div className="absolute inset-0">
        <motion.div
          className="absolute top-0 left-1/4 w-[600px] h-[600px] bg-gradient-to-br from-primary/5 to-transparent rounded-full blur-3xl"
          animate={{
            scale: [1, 1.2, 1],
            rotate: [0, 90, 0],
          }}
          transition={{
            duration: 20,
            repeat: Infinity,
            ease: 'linear',
          }}
        />
        <motion.div
          className="absolute bottom-0 right-1/4 w-[600px] h-[600px] bg-gradient-to-br from-accent/5 to-transparent rounded-full blur-3xl"
          animate={{
            scale: [1.2, 1, 1.2],
            rotate: [90, 0, 90],
          }}
          transition={{
            duration: 20,
            repeat: Infinity,
            ease: 'linear',
          }}
        />
      </div>

      <div className="relative max-w-7xl mx-auto px-6">
        <motion.div
          className="text-center mb-20"
          style={{ scale, opacity }}
        >
          <motion.span
            className="inline-block px-4 py-2 bg-primary/10 text-primary rounded-full text-sm font-semibold mb-6"
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
          >
            技术栈
          </motion.span>

          <h2 className="text-4xl lg:text-5xl font-extrabold text-slate-900 mb-6">
            现代技术架构
          </h2>

          <p className="text-xl text-slate-600 max-w-2xl mx-auto">
            采用业界主流技术栈，确保稳定性与可扩展性
          </p>
        </motion.div>

        <div className="grid md:grid-cols-2 gap-8">
          <TechCard tech={techStack.backend} index={0} />
          <TechCard tech={techStack.frontend} index={1} />
        </div>
      </div>
    </section>
  );
}
