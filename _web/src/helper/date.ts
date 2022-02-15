export const dayCount = (d: Date) => {
  const currentDate = new Date();
  return Math.round((d.getTime() - currentDate.getTime()) / (1000*60*60*24));
};
