import { motion, useScroll, useTransform } from 'framer-motion';
import { useRef, useState } from 'react';
import { useShouldReduceMotion } from '../hooks/useIsMobile';

const codeLines = [
  { type: 'comment', content: '# 克隆项目' },
  { type: 'command', content: 'git clone ', value: 'https://github.com/Davied-H/dict-hub.git' },
  { type: 'command', content: 'cd ', value: 'dict-hub' },
  { type: 'empty', content: '' },
  { type: 'comment', content: '# 一键启动' },
  { type: 'command', content: 'chmod ', value: '+x start.sh' },
  { type: 'command', content: './', value: 'start.sh start' },
  { type: 'empty', content: '' },
  { type: 'comment', content: '# 访问应用' },
  { type: 'comment', content: '# 前端: http://localhost:3000' },
  { type: 'comment', content: '# API:  http://localhost:8080' },
];

export default function QuickStart() {
  const containerRef = useRef<HTMLDivElement>(null);
  const [copied, setCopied] = useState(false);
  const reduceMotion = useShouldReduceMotion();

  const { scrollYProgress } = useScroll({
    target: containerRef,
    offset: ['start end', 'center center'],
  });

  const y = useTransform(scrollYProgress, [0, 1], [100, 0]);
  const opacity = useTransform(scrollYProgress, [0, 0.5], [0, 1]);
  const scale = useTransform(scrollYProgress, [0, 0.5], [0.95, 1]);

  const handleCopy = () => {
    const commands = codeLines
      .filter((l) => l.type === 'command')
      .map((l) => l.content + l.value)
      .join('\n');
    navigator.clipboard.writeText(commands);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  // 移动端简化版本
  if (reduceMotion) {
    return (
      <section
        ref={containerRef}
        id="start"
        className="relative py-32 bg-slate-50 overflow-hidden"
      >
        {/* Background pattern */}
        <div className="absolute inset-0 opacity-50">
          <div
            className="absolute inset-0"
            style={{
              backgroundImage: `radial-gradient(circle at 1px 1px, rgb(148 163 184 / 0.2) 1px, transparent 0)`,
              backgroundSize: '40px 40px',
            }}
          />
        </div>

        <div className="relative max-w-4xl mx-auto px-6">
          <div className="text-center mb-16">
            <span className="inline-block px-4 py-2 bg-primary/10 text-primary rounded-full text-sm font-semibold mb-6">
              快速开始
            </span>

            <h2 className="text-4xl lg:text-5xl font-extrabold text-slate-900 mb-6">
              三步启动
            </h2>

            <p className="text-xl text-slate-600">
              只需几分钟，即可在本地运行 Dict-Hub
            </p>
          </div>

          <div className="relative">
            {/* Terminal window */}
            <div className="relative bg-slate-900 rounded-3xl overflow-hidden shadow-2xl">
              {/* Terminal header */}
              <div className="flex items-center justify-between px-6 py-4 bg-slate-800/80 border-b border-slate-700">
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 rounded-full bg-red-500" />
                  <div className="w-3 h-3 rounded-full bg-yellow-500" />
                  <div className="w-3 h-3 rounded-full bg-green-500" />
                </div>
                <span className="text-slate-400 text-sm font-mono">Terminal</span>
                <button
                  onClick={handleCopy}
                  className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-slate-700 hover:bg-slate-600 text-slate-300 text-sm transition-colors"
                >
                  {copied ? (
                    <>
                      <svg className="w-4 h-4 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                      </svg>
                      已复制
                    </>
                  ) : (
                    <>
                      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                      </svg>
                      复制
                    </>
                  )}
                </button>
              </div>

              {/* Code content */}
              <div className="p-8 font-mono text-sm lg:text-base leading-relaxed overflow-x-auto">
                {codeLines.map((line, i) => (
                  <div key={i} className="flex">
                    {line.type === 'comment' && (
                      <span className="text-slate-500">{line.content}</span>
                    )}
                    {line.type === 'command' && (
                      <>
                        <span className="text-cyan-400">{line.content}</span>
                        <span className="text-cyan-200">{line.value}</span>
                      </>
                    )}
                    {line.type === 'empty' && <span>&nbsp;</span>}
                  </div>
                ))}

                {/* Static cursor */}
                <span className="inline-block w-2.5 h-5 bg-slate-400 ml-1 animate-pulse" />
              </div>
            </div>
          </div>
        </div>
      </section>
    );
  }

  // 桌面端完整动画版本
  return (
    <section
      ref={containerRef}
      id="start"
      className="relative py-32 bg-slate-50 overflow-hidden"
    >
      {/* Background pattern */}
      <div className="absolute inset-0 opacity-50">
        <div
          className="absolute inset-0"
          style={{
            backgroundImage: `radial-gradient(circle at 1px 1px, rgb(148 163 184 / 0.2) 1px, transparent 0)`,
            backgroundSize: '40px 40px',
          }}
        />
      </div>

      <div className="relative max-w-4xl mx-auto px-6">
        <motion.div
          className="text-center mb-16"
          style={{ y, opacity }}
        >
          <motion.span
            className="inline-block px-4 py-2 bg-primary/10 text-primary rounded-full text-sm font-semibold mb-6"
            initial={{ opacity: 0, scale: 0.8 }}
            whileInView={{ opacity: 1, scale: 1 }}
            viewport={{ once: true }}
          >
            快速开始
          </motion.span>

          <h2 className="text-4xl lg:text-5xl font-extrabold text-slate-900 mb-6">
            三步启动
          </h2>

          <p className="text-xl text-slate-600">
            只需几分钟，即可在本地运行 Dict-Hub
          </p>
        </motion.div>

        <motion.div
          style={{ scale, opacity }}
          className="relative"
        >
          {/* Terminal window */}
          <motion.div
            className="relative bg-slate-900 rounded-3xl overflow-hidden shadow-2xl"
            initial={{ opacity: 0, y: 40 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.8, ease: [0.16, 1, 0.3, 1] }}
          >
            {/* Terminal header */}
            <div className="flex items-center justify-between px-6 py-4 bg-slate-800/80 border-b border-slate-700">
              <div className="flex items-center gap-2">
                <div className="w-3 h-3 rounded-full bg-red-500" />
                <div className="w-3 h-3 rounded-full bg-yellow-500" />
                <div className="w-3 h-3 rounded-full bg-green-500" />
              </div>
              <span className="text-slate-400 text-sm font-mono">Terminal</span>
              <motion.button
                onClick={handleCopy}
                className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-slate-700 hover:bg-slate-600 text-slate-300 text-sm transition-colors"
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                {copied ? (
                  <>
                    <svg className="w-4 h-4 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    已复制
                  </>
                ) : (
                  <>
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                    </svg>
                    复制
                  </>
                )}
              </motion.button>
            </div>

            {/* Code content */}
            <div className="p-8 font-mono text-sm lg:text-base leading-relaxed overflow-x-auto">
              {codeLines.map((line, i) => (
                <motion.div
                  key={i}
                  className="flex"
                  initial={{ opacity: 0, x: -20 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: i * 0.05 }}
                >
                  {line.type === 'comment' && (
                    <span className="text-slate-500">{line.content}</span>
                  )}
                  {line.type === 'command' && (
                    <>
                      <span className="text-cyan-400">{line.content}</span>
                      <span className="text-cyan-200">{line.value}</span>
                    </>
                  )}
                  {line.type === 'empty' && <span>&nbsp;</span>}
                </motion.div>
              ))}

              {/* Cursor blink */}
              <motion.span
                className="inline-block w-2.5 h-5 bg-slate-400 ml-1"
                animate={{ opacity: [1, 0] }}
                transition={{ duration: 0.8, repeat: Infinity }}
              />
            </div>
          </motion.div>

          {/* Decorative elements */}
          <motion.div
            className="absolute -top-6 -right-6 w-24 h-24 bg-gradient-to-br from-primary to-accent rounded-2xl opacity-20 blur-xl"
            animate={{
              rotate: [0, 90, 0],
              scale: [1, 1.2, 1],
            }}
            transition={{ duration: 8, repeat: Infinity }}
          />
          <motion.div
            className="absolute -bottom-6 -left-6 w-32 h-32 bg-gradient-to-br from-accent to-primary rounded-full opacity-20 blur-xl"
            animate={{
              rotate: [0, -90, 0],
              scale: [1.2, 1, 1.2],
            }}
            transition={{ duration: 10, repeat: Infinity }}
          />
        </motion.div>
      </div>
    </section>
  );
}
