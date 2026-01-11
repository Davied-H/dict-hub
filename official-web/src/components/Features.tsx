import { motion, useScroll, useTransform, useInView } from 'framer-motion';
import { useRef } from 'react';
import { useShouldReduceMotion } from '../hooks/useIsMobile';

const features = [
  {
    icon: '🔍',
    title: '快速查词',
    description: '基于高效索引的毫秒级词条检索，支持模糊搜索与智能建议',
    color: 'from-blue-500 to-cyan-500',
    shadowColor: 'shadow-blue-500/20',
    demo: (
      <div className="space-y-2">
        <div className="flex items-center gap-2 px-3 py-2 bg-white/80 rounded-lg text-sm">
          <span className="text-slate-400">搜索:</span>
          <span className="text-slate-700 font-medium">elo</span>
          <span className="w-0.5 h-4 bg-blue-500 animate-pulse" />
        </div>
        <div className="text-xs text-slate-500 pl-2">
          <div className="flex items-center gap-2 py-1">
            <span className="w-1.5 h-1.5 rounded-full bg-blue-500" />
            eloquent
          </div>
          <div className="flex items-center gap-2 py-1 opacity-60">
            <span className="w-1.5 h-1.5 rounded-full bg-slate-300" />
            eloquence
          </div>
        </div>
      </div>
    ),
  },
  {
    icon: '📚',
    title: '多词典支持',
    description: '支持加载多个 MDX 格式词典，灵活切换与管理',
    color: 'from-violet-500 to-purple-500',
    shadowColor: 'shadow-violet-500/20',
    demo: (
      <div className="space-y-2">
        {['牛津高阶', '朗文当代', '柯林斯'].map((dict, i) => (
          <div
            key={dict}
            className={`flex items-center gap-2 px-3 py-1.5 rounded-lg text-xs ${
              i === 0 ? 'bg-violet-100 text-violet-700 font-medium' : 'bg-white/60 text-slate-500'
            }`}
          >
            <span className={`w-2 h-2 rounded-full ${i === 0 ? 'bg-violet-500' : 'bg-slate-300'}`} />
            {dict}
          </div>
        ))}
      </div>
    ),
  },
  {
    icon: '📊',
    title: '词频统计',
    description: '自动记录查词频率，帮助了解学习进度与重点词汇',
    color: 'from-emerald-500 to-teal-500',
    shadowColor: 'shadow-emerald-500/20',
    demo: (
      <div className="space-y-2">
        {[
          { word: 'eloquent', count: 12, width: '100%' },
          { word: 'profound', count: 8, width: '66%' },
          { word: 'aesthetic', count: 5, width: '42%' },
        ].map((item) => (
          <div key={item.word} className="space-y-1">
            <div className="flex justify-between text-xs">
              <span className="text-slate-600">{item.word}</span>
              <span className="text-emerald-600 font-medium">{item.count}次</span>
            </div>
            <div className="h-1.5 bg-slate-100 rounded-full overflow-hidden">
              <motion.div
                className="h-full bg-gradient-to-r from-emerald-500 to-teal-500 rounded-full"
                initial={{ width: 0 }}
                whileInView={{ width: item.width }}
                viewport={{ once: true }}
                transition={{ duration: 1, delay: 0.5 }}
              />
            </div>
          </div>
        ))}
      </div>
    ),
  },
  {
    icon: '📝',
    title: '历史记录',
    description: '完整的查询历史追踪，支持回顾与导出',
    color: 'from-amber-500 to-orange-500',
    shadowColor: 'shadow-amber-500/20',
    demo: (
      <div className="space-y-1.5">
        {[
          { word: 'serendipity', time: '刚刚' },
          { word: 'ephemeral', time: '2分钟前' },
          { word: 'ubiquitous', time: '5分钟前' },
        ].map((item, i) => (
          <div
            key={item.word}
            className="flex items-center justify-between px-3 py-1.5 bg-white/60 rounded-lg text-xs"
            style={{ opacity: 1 - i * 0.2 }}
          >
            <span className="text-slate-700">{item.word}</span>
            <span className="text-slate-400">{item.time}</span>
          </div>
        ))}
      </div>
    ),
  },
  {
    icon: '🎨',
    title: '现代界面',
    description: '简洁优雅的用户界面，支持深色模式切换',
    color: 'from-pink-500 to-rose-500',
    shadowColor: 'shadow-pink-500/20',
    demo: (
      <div className="flex gap-2 items-center justify-center">
        <div className="w-12 h-8 bg-white rounded-lg border border-slate-200 flex items-center justify-center">
          <div className="w-6 h-4 bg-slate-100 rounded" />
        </div>
        <div className="text-slate-300">⇄</div>
        <div className="w-12 h-8 bg-slate-800 rounded-lg flex items-center justify-center">
          <div className="w-6 h-4 bg-slate-700 rounded" />
        </div>
      </div>
    ),
  },
  {
    icon: '🚀',
    title: '高性能架构',
    description: 'Go 后端 + React 前端，极致性能与用户体验',
    color: 'from-red-500 to-rose-600',
    shadowColor: 'shadow-red-500/20',
    demo: (
      <div className="space-y-2">
        <div className="flex items-center gap-2 text-xs">
          <span className="text-slate-500">响应时间</span>
          <div className="flex-1 h-1.5 bg-slate-100 rounded-full overflow-hidden">
            <motion.div
              className="h-full bg-gradient-to-r from-red-500 to-rose-500 rounded-full"
              initial={{ width: 0 }}
              whileInView={{ width: '15%' }}
              viewport={{ once: true }}
              transition={{ duration: 0.8, delay: 0.3 }}
            />
          </div>
          <span className="text-red-600 font-bold">12ms</span>
        </div>
        <div className="flex justify-center gap-4 text-xs text-slate-500">
          <span className="flex items-center gap-1">
            <span className="w-2 h-2 rounded-full bg-cyan-500" />
            Go
          </span>
          <span className="flex items-center gap-1">
            <span className="w-2 h-2 rounded-full bg-blue-500" />
            React
          </span>
        </div>
      </div>
    ),
  },
];

function FeatureCard({
  feature,
  index,
  reduceMotion,
}: {
  feature: (typeof features)[0];
  index: number;
  reduceMotion: boolean;
}) {
  const cardRef = useRef<HTMLDivElement>(null);
  const isInView = useInView(cardRef, { once: true, margin: '-100px' });

  // 移动端使用简化的动画
  if (reduceMotion) {
    return (
      <div
        ref={cardRef}
        className={`group relative bg-white rounded-3xl border border-slate-200/80 overflow-hidden hover:shadow-2xl ${feature.shadowColor} transition-shadow duration-500`}
      >
        <div className={`h-1.5 bg-gradient-to-r ${feature.color}`} />
        <div className="p-8">
          <div className="relative mb-6">
            <div
              className={`relative w-16 h-16 bg-gradient-to-br ${feature.color} rounded-2xl flex items-center justify-center text-3xl shadow-lg ${feature.shadowColor}`}
            >
              {feature.icon}
            </div>
          </div>
          <h3 className="text-xl font-bold text-slate-900 mb-3">
            {feature.title}
          </h3>
          <p className="text-slate-600 leading-relaxed mb-6">
            {feature.description}
          </p>
          <div className="p-4 bg-gradient-to-br from-slate-50 to-slate-100/50 rounded-2xl border border-slate-100">
            {feature.demo}
          </div>
        </div>
      </div>
    );
  }

  return (
    <motion.div
      ref={cardRef}
      initial={{ opacity: 0, y: 80 }}
      animate={isInView ? { opacity: 1, y: 0 } : {}}
      transition={{
        duration: 0.7,
        delay: index * 0.1,
        ease: [0.16, 1, 0.3, 1],
      }}
      whileHover={{ y: -12, transition: { duration: 0.3 } }}
      className={`group relative bg-white rounded-3xl border border-slate-200/80 overflow-hidden hover:shadow-2xl ${feature.shadowColor} transition-shadow duration-500`}
    >
      {/* Top gradient bar */}
      <div className={`h-1.5 bg-gradient-to-r ${feature.color}`} />

      <div className="p-8">
        {/* Icon with animated background */}
        <div className="relative mb-6">
          <motion.div
            className={`absolute inset-0 bg-gradient-to-r ${feature.color} rounded-2xl blur-xl opacity-0 group-hover:opacity-30`}
            transition={{ duration: 0.5 }}
          />
          <motion.div
            className={`relative w-16 h-16 bg-gradient-to-br ${feature.color} rounded-2xl flex items-center justify-center text-3xl shadow-lg ${feature.shadowColor}`}
            whileHover={{
              scale: 1.1,
              rotate: [0, -10, 10, -5, 0],
            }}
            transition={{ duration: 0.5 }}
          >
            {feature.icon}
          </motion.div>
        </div>

        {/* Title */}
        <h3 className="text-xl font-bold text-slate-900 mb-3 group-hover:text-transparent group-hover:bg-clip-text group-hover:bg-gradient-to-r group-hover:from-slate-900 group-hover:to-slate-600 transition-all duration-300">
          {feature.title}
        </h3>

        {/* Description */}
        <p className="text-slate-600 leading-relaxed mb-6">
          {feature.description}
        </p>

        {/* Interactive demo area */}
        <motion.div
          className="p-4 bg-gradient-to-br from-slate-50 to-slate-100/50 rounded-2xl border border-slate-100"
          initial={{ opacity: 0.5, scale: 0.95 }}
          whileInView={{ opacity: 1, scale: 1 }}
          viewport={{ once: true }}
          transition={{ delay: 0.3 + index * 0.1 }}
        >
          {feature.demo}
        </motion.div>
      </div>

      {/* Hover corner decoration */}
      <div
        className={`absolute -bottom-20 -right-20 w-40 h-40 bg-gradient-to-br ${feature.color} rounded-full opacity-0 group-hover:opacity-10 blur-2xl transition-opacity duration-500`}
      />
    </motion.div>
  );
}

export default function Features() {
  const containerRef = useRef<HTMLDivElement>(null);
  const reduceMotion = useShouldReduceMotion();

  const { scrollYProgress } = useScroll({
    target: containerRef,
    offset: ['start end', 'center center'],
  });

  const titleY = useTransform(scrollYProgress, [0, 1], [60, 0]);
  const titleOpacity = useTransform(scrollYProgress, [0, 0.5], [0, 1]);

  return (
    <section
      ref={containerRef}
      id="features"
      className="relative py-32 bg-gradient-to-b from-white via-slate-50 to-white overflow-hidden"
    >
      {/* Animated background elements - 仅在桌面端显示 */}
      {!reduceMotion && (
        <div className="absolute inset-0 pointer-events-none overflow-hidden">
          <motion.div
            className="absolute top-20 left-10 w-72 h-72 bg-blue-500/5 rounded-full blur-3xl"
            animate={{
              x: [0, 50, 0],
              y: [0, 30, 0],
            }}
            transition={{ duration: 15, repeat: Infinity, ease: 'easeInOut' }}
          />
          <motion.div
            className="absolute bottom-20 right-10 w-96 h-96 bg-purple-500/5 rounded-full blur-3xl"
            animate={{
              x: [0, -30, 0],
              y: [0, -50, 0],
            }}
            transition={{ duration: 18, repeat: Infinity, ease: 'easeInOut' }}
          />
          <motion.div
            className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-gradient-to-r from-cyan-500/5 to-blue-500/5 rounded-full blur-3xl"
            animate={{ rotate: 360 }}
            transition={{ duration: 60, repeat: Infinity, ease: 'linear' }}
          />
        </div>
      )}

      <div className="relative max-w-7xl mx-auto px-6">
        {/* Section header */}
        {reduceMotion ? (
          <div className="text-center mb-20">
            <div className="inline-flex items-center gap-2 px-5 py-2.5 bg-gradient-to-r from-blue-500/10 to-cyan-500/10 rounded-full mb-6">
              <span className="relative flex h-2 w-2">
                <span className="relative inline-flex rounded-full h-2 w-2 bg-blue-500" />
              </span>
              <span className="text-primary font-semibold text-sm">功能特性</span>
            </div>
            <h2 className="text-4xl lg:text-6xl font-extrabold mb-6">
              <span className="text-slate-900">为</span>
              <span className="text-gradient">效率</span>
              <span className="text-slate-900">而生</span>
            </h2>
            <p className="text-xl text-slate-600 max-w-2xl mx-auto">
              Dict-Hub 提供一系列强大功能，让词汇学习与查询更加高效便捷
            </p>
          </div>
        ) : (
          <motion.div
            className="text-center mb-20"
            style={{ y: titleY, opacity: titleOpacity }}
          >
            <motion.div
              className="inline-flex items-center gap-2 px-5 py-2.5 bg-gradient-to-r from-blue-500/10 to-cyan-500/10 rounded-full mb-6"
              initial={{ opacity: 0, scale: 0.8 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5 }}
            >
              <span className="relative flex h-2 w-2">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75" />
                <span className="relative inline-flex rounded-full h-2 w-2 bg-blue-500" />
              </span>
              <span className="text-primary font-semibold text-sm">功能特性</span>
            </motion.div>

            <motion.h2
              className="text-4xl lg:text-6xl font-extrabold mb-6"
              initial={{ opacity: 0, y: 30 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.6, delay: 0.1 }}
            >
              <span className="text-slate-900">为</span>
              <span className="text-gradient">效率</span>
              <span className="text-slate-900">而生</span>
            </motion.h2>

            <motion.p
              className="text-xl text-slate-600 max-w-2xl mx-auto"
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.6, delay: 0.2 }}
            >
              Dict-Hub 提供一系列强大功能，让词汇学习与查询更加高效便捷
            </motion.p>
          </motion.div>
        )}

        {/* Feature cards grid */}
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          {features.map((feature, i) => (
            <FeatureCard key={feature.title} feature={feature} index={i} reduceMotion={reduceMotion} />
          ))}
        </div>
      </div>
    </section>
  );
}
